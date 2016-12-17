package warpify

import "github.com/pressly/warpdrive/config"
import "github.com/pressly/warpdrive/lib/folder"
import "github.com/blang/semver"
import "github.com/pressly/warpdrive/data"

const (
	_ EventKind = 0

	// Err if something goes wrong, Err event will be sent
	Err
	// NoUpdate means there is no update available at the moment
	NoUpdate
	// UpdateAvailable there is an update available for download
	UpdateAvailable
	// UpdateDownloading downloading has started
	UpdateDownloading
	// UpdateDownloaded downlaoing has completed
	// at this moment, a callback from objective c or java should restart the app
	UpdateDownloaded
	// HardUpdateAvailable will be called when production release
	// has hard available and does not have soft update
	HardUpdateAvailable
)

// Setup we need to setup the app
func Setup(bundleVersion, bundlePath, documentPath, platform, productionName string, automaticUpdate bool) {
	conf.bundleVersion = bundleVersion
	conf.bundlePath = bundlePath
	conf.documentPath = bundlePath
	conf.productionName = productionName
	conf.automaticUpdate = automaticUpdate
	conf.platform = platform

	conf.pubSub = newPubSub()
}

// SourcePath returns the proper path for react-native app to start the process
func SourcePath() string {
	return ""
}

// Start accepts a callback and start the process of checking the version
// and download and restart the
func Start(callback Callback) {
	// attach the following events to given callback
	subscribe(Err, callback)
	subscribe(NoUpdate, callback)
	subscribe(UpdateAvailable, callback)
	subscribe(UpdateDownloading, callback)
	subscribe(UpdateDownloaded, callback)

	go func() {
		err := Process()
		if err != nil {
			conf.pubSub.Publish(createEvent(Err, err.Error()))
		}
	}()
}

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
	api := makeApi(warpFile, versionMap)

	var autoUpdateRelease map[string]*data.Release
	var releases []map[string]*data.Release
	// we need to loop through all the available configs
	for _, cycleConfig := range warpFile.Cycles {
		release, err := api.checkVersion(appID, cycleConfig.ID, conf.platform, versionMap.CurrentVersion(cycleConfig.Name))
		if err != nil {
			conf.pubSub.Publish(createEvent(Err, err.Error()))
		} else {
			if conf.productionName == cycleConfig.Name {
				autoUpdateRelease = release
			} else {
				releases = append(releases, release)
			}
		}
	}

	if autoUpdateRelease != nil {
		//we need to check whether update has soft key
		softRelease, ok := autoUpdateRelease["soft"]
		if ok {
			conf.pubSub.Publish(createEvent(UpdateDownloading, softRelease))
			r, err := api.downloadVersion(appID, softRelease.CycleID, softRelease.ID)
			if err != nil {
				conf.pubSub.Publish(createEvent(Err, err.Error()))
			} else {
				if r != nil {

				}
			}
		}
	}

	if len(releases) > 0 {
		//conf.pubSub()
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
			cycleName = conf.productionName
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
func subscribe(eventKind EventKind, callback Callback) {
	conf.pubSub.Subscribe(eventKind, callback)
}

// unsubscribe as it stands, it unsubscribes the any associate
// callback to specific event type. Mainly it's being used for clean up
func unsubscribe(eventKind EventKind) {
	conf.pubSub.Unsubscribe(eventKind)
}
