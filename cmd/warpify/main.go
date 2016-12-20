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
func Setup(bundleVersion, bundlePath, documentPath, platform, defaultCycle string, forceUpdate bool) error {
	conf.bundleVersion = bundleVersion
	conf.bundlePath = bundlePath
	conf.documentPath = bundlePath
	conf.defaultCycle = defaultCycle
	conf.forceUpdate = forceUpdate
	conf.platform = platform

	conf.pubSub = newPubSub()

	// load versionMap and warpFile
	warpFile, err := getWarpFile()
	if err != nil {
		return err
	}

	conf.warpFile = warpFile

	// we need to create api
	conf.api = makeAPI(warpFile)

	return nil
}

// DownloadRelease it starts download and save the bundle inside path
func DownloadRelease(cycleID, releaseID int64, version string) error {
	appID := conf.warpFile.App.ID

	r, err := conf.api.downloadVersion(appID, cycleID, releaseID)
	if err != nil {
		return err
	}

	path := warpBundlePath(appID, cycleID, version)

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
// if there is no downlaod available, simply return an empty string
// on native side, it will be replaced by default value
// also this method must be called after Setup function is called with no errors
func SourcePath() string {
	versionMap, err := getVersionMap()
	if err != nil {
		return ""
	}

	// So we grab the ActiveCycle to find the current active cycle
	// then we look for current version
	// if current version is bundle then, we return "" and let the
	// native code handle the URL
	// if not then we need to generate a path to new value

	version := versionMap.Version(versionMap.ActiveCycle)

	// we need to search for current version inside avaiable

	isBundle, ok := version.Available[version.Current]
	if !ok || isBundle {
		return ""
	}

	appID := conf.warpFile.App.ID

	cycleConfig, err := conf.warpFile.GetCycle(versionMap.ActiveCycle)
	if err != nil {
		return ""
	}

	// we need to point to main.jsbundle file
	path := filepath.Join(warpBundlePath(appID, cycleConfig.ID, version.Current), "main.jsbundle")

	// start the process but only for defaultCyle
	process(true)

	return path
}

// Start accepts a callback and start the process of checking the version
// and download and restart the
func Start(callback Callback) {
	if callback != nil {
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
	}

	go func() {
		err := process(callback == nil)
		if err != nil {
			conf.pubSub.Publish(Err, err.Error())
		}
	}()
}

// process is the main logic to handle all the updates
func process(justDefaultCycle bool) error {
	versionMap, err := getVersionMap()
	if err != nil {
		return err
	}

	appID := conf.warpFile.App.ID

	var releases []map[string]*data.Release

	// first we need to see if there is a version for deafultCycle
	// if there is one, then there is no reason to let app know about other version
	defaultCycleConfig, err := conf.warpFile.GetCycle(conf.defaultCycle)
	if err != nil {
		// we are terminating the process because defaultConfig is not found
		return err
	}

	// lets check if default cycle has a new update
	currentVersion := versionMap.CurrentVersion(defaultCycleConfig.Name)
	releaseMap, err := conf.api.checkVersion(appID, defaultCycleConfig.ID, conf.platform, currentVersion)
	if err != nil {
		// at this point we are not terminate the process, we simple send an event
		// to notify the default cycle has some issue
		conf.pubSub.Publish(Err, err.Error())
	} else {
		// we need to check if soft update is available

		if softRelease, ok := releaseMap["soft"]; ok {
			// we need to check if forceUpdate is enabled
			// if forceUpdate is enable, then we simple download the update and update the version
			// and we should call the objective-c and java part for restart the app
			if conf.forceUpdate {
				err = DownloadRelease(defaultCycleConfig.ID, softRelease.ID, currentVersion)
				if err != nil {
					return err
				}

				// We need to call the native to restart the app
				Reload(warpBundlePath(appID, defaultCycleConfig.ID, currentVersion))
				return nil
			}

			// since force update is not enabled, then we are adding the releases
			// and let client decides whether app needs to be updated or not
			releases = append(releases, releaseMap)
		}
	}

	if justDefaultCycle {
		return nil
	}

	// we need to loop through all available configs
	for _, cycleConfig := range conf.warpFile.Cycles {
		// we don't need to check the default cycle again
		if cycleConfig.Name == conf.defaultCycle {
			continue
		}

		releaseMap, err := conf.api.checkVersion(appID, cycleConfig.ID, conf.platform, versionMap.CurrentVersion(cycleConfig.Name))
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
	if exists, _ := folder.PathExists(path); exists {
		err := versionMap.Load(conf.documentPath)
		if err != nil {
			return nil, err
		}
	} else {
		// this is the first time we are creating versions.warp
		// this file will be saved in documentPath/warpdrive/versions.warp

		version, err := semver.Make(conf.bundleVersion)
		if err != nil {
			return nil, err
		}

		// if the bundle contains pre value then it consider
		// default cycle name, make sure in your code, defaultCycle and
		// pre in bundle version was the same
		var cycleName string
		if len(version.Pre) > 0 {
			cycleName = version.Pre[0].String()
		} else {
			cycleName = conf.defaultCycle
		}

		// set the first version inside versions.warp
		// and because the first version points in bundle, then
		// isBundle in `SetCurrentVersion` will be true
		versionMap.ActiveCycle = cycleName
		versionMap.SetCurrentVersion(cycleName, conf.bundleVersion, true, false)

		// save the file on disk
		err = versionMap.Save(path)
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

func releaseToString(releases []map[string]*data.Release) string {
	str, _ := json.Marshal(releases)
	return string(str)
}

func warpBundlePath(appID, cycleID int64, version string) string {
	path := fmt.Sprintf("warpdrive/warp.%d.%d.%s", appID, cycleID, version)
	return filepath.Join(conf.documentPath, path)
}
