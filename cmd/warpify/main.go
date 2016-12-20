package warpify

import (
	"fmt"
	"os"
	"path/filepath"

	"encoding/json"

	"github.com/blang/semver"
	"github.com/pressly/warpdrive/config"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/folder"
	"github.com/pressly/warpdrive/lib/warp"
)

const (
	_ int = 0

	// Err if something goes wrong, Err event will be sent
	Err
	// NoUpdate happens when there is no update avaiable to download
	NoUpdate
	// Available means an update is avaliable but has not been downlaoded
	// it requires user request for downloading the update
	Available
	// Downloading means update is being downloaded
	Downloading
	// Downloaded means requested update has been downloaded
	Downloaded
)

func releaseToString(releases []map[string]*data.Release) string {
	str, _ := json.Marshal(releases)
	return string(str)
}

// SetReload assing an func call from native target for reloading
// react-native
func SetReload(reloadTask Callback) {
	conf.reloadTask = reloadTask
}

// Reload passes path of assets and force the react-native to reload.
// Target native code must register a proepr function for rreloading the
// react-native or this method returns an error
func Reload(path string) {
	conf.reloadTask.Do(0, path)
}

// Setup we need to setup the app
func Setup(bundleVersion, bundlePath, documentPath, platform, defaultCycle string, forceUpdate bool) {
	conf.bundleVersion = bundleVersion
	conf.bundlePath = bundlePath
	conf.documentPath = bundlePath
	conf.defaultCycle = defaultCycle
	conf.forceUpdate = forceUpdate
	conf.platform = platform

	conf.pubSub = newPubSub()
}

func warpBundlePath(appID, cycleID, releaseID int64) string {
	path := fmt.Sprintf("warpdrive/warp.%d.%d.%d", appID, cycleID, releaseID)
	return filepath.Join(conf.documentPath, path)
}

// DownloadRelease it starts download and save the bundle inside path
func DownloadRelease(api *api, appID, cycleID, releaseID int64) error {
	r, err := api.downloadVersion(appID, cycleID, releaseID)
	if err != nil {
		return err
	}

	path := warpBundlePath(appID, cycleID, releaseID)

	// need to make fodler
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// extract the files inside that path
	err = warp.Extract(r, path)
	if err != nil {
		return err
	}

	return nil
}

// SourcePath returns the proper path for react-native app to start the process
func SourcePath() string {
	return ""
}

// Start accepts a callback and start the process of checking the version
// and download and restart the
func Start(callback Callback) {
	// we are going to detach all calbacks first
	unsubscribe(Err)
	unsubscribe(NoUpdate)
	unsubscribe(Available)
	unsubscribe(Downloading)
	unsubscribe(Downloaded)

	// attach the following events to given callback
	subscribe(Err, callback)
	subscribe(NoUpdate, callback)
	subscribe(Available, callback)
	subscribe(Downloading, callback)
	subscribe(Downloaded, callback)

	go func() {
		err := Process()
		if err != nil {
			conf.pubSub.Publish(Err, err.Error())
		}
	}()
}

// Process is the main logic to handle all the updates
func Process() error {
	// load versionMap and warpFile
	warpFile, err := getWarpFile()
	if err != nil {
		return err
	}

	versionMap, err := getVersionMap()
	if err != nil {
		return err
	}

	appID := warpFile.App.ID

	// we need to create api
	api := makeAPI(warpFile)

	var releases []map[string]*data.Release

	// first we need to see if there is a version for deafultCycle
	// if there is one, then there is no reason to let app know about other version
	defaultCycleConfig, err := warpFile.GetCycle(conf.defaultCycle)
	if err != nil {
		// we are terminating the process because defaultConfig is not found
		return err
	}

	// lets check if default cycle has a new update
	releaseMap, err := api.checkVersion(appID, defaultCycleConfig.ID, conf.platform, versionMap.CurrentVersion(defaultCycleConfig.Name))
	if err != nil {
		// at this point we are not terminate the process, we simple send an event
		// to notify the default cycle has some issue
		conf.pubSub.Publish(Err, err.Error())
	} else {
		// we need to check if soft update is available
		softRelease, ok := releaseMap["soft"]
		if ok {
			// we need to check if forceUpdate is enabled
			// if forceUpdate is enable, then we simple download the update and update the version
			// and we should call the objective-c and java part for restart the app
			if conf.forceUpdate {
				err = DownloadRelease(api, appID, defaultCycleConfig.ID, softRelease.ID)
				if err != nil {
					return err
				}

				// We need to call the native to restart the app
				Reload(warpBundlePath(appID, defaultCycleConfig.ID, softRelease.ID))
				return nil
			}
		}
	}

	// we need to loop through all available configs
	for _, cycleConfig := range warpFile.Cycles {
		// we don't need to check the default cycle again
		if cycleConfig.Name == conf.defaultCycle {
			continue
		}

		releaseMap, err := api.checkVersion(appID, cycleConfig.ID, conf.platform, versionMap.CurrentVersion(cycleConfig.Name))
		if err != nil {
			conf.pubSub.Publish(Err, err.Error())
		} else {
			releases = append(releases, releaseMap)
		}
	}

	if len(releases) > 0 {
		conf.pubSub.Publish(Available, releaseToString(releases))
	} else {
		conf.pubSub.Publish(NoUpdate, "")
	}

	return nil
}

// warpFile loads the WarpFile from bundle path
func getWarpFile() (*config.ClientConfig, error) {
	clientConfig := config.NewClientConfigsForMobile(conf.bundlePath)

	err := clientConfig.Load()
	if err != nil {
		return nil, err
	}

	return clientConfig, nil
}

// versionMap loads versions.warp from documentPath, if it exists,
// if not we created a new versions.warp and save it to document folder
func getVersionMap() (*config.VersionMap, error) {
	var versionMap config.VersionMap

	path := config.VersionPath(conf.documentPath)
	exists, _ := folder.PathExists(path)
	if exists {
		err := versionMap.Load(conf.documentPath)
		if err != nil {
			return nil, err
		}
	} else {
		version, err := semver.Make(conf.bundleVersion)
		if err != nil {
			return nil, err
		}

		var cycleName string
		if len(version.Pre) > 0 {
			cycleName = version.Pre[0].String()
		} else {
			cycleName = conf.defaultCycle
		}

		versionMap.SetCurrentVersion(cycleName, conf.bundleVersion, false)
		err = versionMap.Save(conf.documentPath)
		if err != nil {
			return nil, err
		}
	}

	return &versionMap, nil
}

// subscribe this is a easy to use method to expose to objective-c and jave
// so they can bind their callbacks to known EventKinds
func subscribe(eventKind int, callback Callback) {
	conf.pubSub.Subscribe(eventKind, callback)
}

// unsubscribe as it stands, it unsubscribes the any associate
// callback to specific event type. Mainly it's being used for clean up
func unsubscribe(eventKind int) {
	conf.pubSub.Unsubscribe(eventKind)
}
