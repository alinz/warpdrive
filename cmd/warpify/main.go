package warpify

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"encoding/json"

	"io"

	"strings"

	"github.com/blang/semver"
	"github.com/pressly/warpdrive/config"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/folder"
	"github.com/pressly/warpdrive/lib/warp"
)

// Callback is an interface for calling native with some values
type Callback interface {
	Do(kind int, value string)
}

// Cycles returns the stringify of list of available cycles in the bundles
// e.g. [{ id: 1, name: "dev" }, ...]
func Cycles() (string, error) {
	var results []map[string]interface{}

	for _, cycle := range conf.warpFile.Cycles {
		cycleInfo := make(map[string]interface{})
		cycleInfo["id"] = cycle.ID
		cycleInfo["name"] = cycle.Name

		results = append(results, cycleInfo)
	}

	b, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// LocalVersions we can not returns an array of data.Release because gomobile
// won't generate proepr objective-c binding, what we do here is serliazie the
// array into string and return the value. The string will be parsed in javascript side
func LocalVersions(cycleID int64) (string, error) {
	warpPath := filepath.Join(conf.documentPath, "warpdrive")

	folders, err := folder.ListFolders(warpPath)
	if err != nil {
		return "", err
	}

	var results []map[string]string

	for _, folder := range folders {
		//warp.<appID>.<cycleID>.<version>
		segments := strings.Split(folder, ".")
		if len(segments) == 4 {
			targetCycleID, err := strconv.ParseInt(segments[2], 10, 64)
			if err != nil {
				continue
			}

			if targetCycleID == cycleID {
				versionMap := make(map[string]string)
				versionMap["version"] = segments[3]
				results = append(results, versionMap)
			}
		}
	}

	b, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// RemoteVersions we can not returns an array of data.Release because gomobile
// won't generate proepr objective-c binding, what we do here is serliazie the
// array into string and return the value. The string will be parsed in javascript side
func RemoteVersions(cycleID int64) (string, error) {
	appID := conf.warpFile.App.ID

	releases, err := conf.api.remoteVersions(appID, cycleID)
	if err != nil {
		return "", err
	}

	var results []map[string]interface{}
	for _, release := range releases {
		releaseMap := make(map[string]interface{})
		releaseMap["version"] = release.Version
		releaseMap["note"] = release.Note

		results = append(results, releaseMap)
	}

	b, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// SetReload assing an func call from native target for reloading
// react-native, this is only need to be used if forceUpdate is enabled.
func SetReload(reloadTask Callback) {
	conf.reloadTask = reloadTask

	// once the reloadTask is set, we are going to close the reloadReady channel
	// which triggers the go to go ahead and run the task.
	close(conf.reloadReady)
}

// Reload accepts a version. The version needs to be already downloaded or
// an error will be return. version also needs to follow semantic versioning
// the preRelease section of the version targets the cycle's name
func Reload(cycleID int64, versionStr string) error {
	version, err := semver.Make(versionStr)
	if err != nil {
		return err
	}

	// we need to find cycle information about this version
	cycleName := cycleNameFromVersion(&version)
	cycleConfig, err := conf.warpFile.GetCycle(cycleName)
	if err != nil {
		return err
	}

	// before we update the the path, we need to make sure
	// the versionMap exists
	versionMap, err := getVersionMap()
	if err != nil {
		return err
	}

	versionMap.SetActiveCycle(cycleConfig.Name)
	// passing fales to isBundle argument is safe. becuase internally we are checking if the version
	// is indeed exists in bundle folder or document.
	isSourceInBundle := versionMap.SetCurrentVersion(cycleConfig.Name, version.String(), false)
	versionMap.Save(conf.documentPath)

	var path string

	if isSourceInBundle {
		path = conf.bundlePath
	} else {
		path = warpBundlePath(conf.warpFile.App.ID, cycleConfig.ID, version.String())
	}

	// NOTE: one more thing here, since, we can't get the bundle path for android,
	// the path for bundle file in android will be wrong, what we have to do
	// in android side is we need to check the path with temporary bundle and if the are match
	// we need to return null. Null in android triggers internal call which automatically
	// loads the bundle file. I don't want to change this here. The go code must be completely
	// agnostic to platform.
	conf.reloadTask.Do(0, filepath.Join(path, "main.jsbundle"))

	return nil
}

// Setup we need to setup the app
func Setup(bundleVersion, bundlePath, documentPath, platform, defaultCycle string, forceUpdate bool) error {
	conf.bundleVersion = bundleVersion
	conf.bundlePath = bundlePath
	conf.documentPath = documentPath
	conf.defaultCycle = defaultCycle
	conf.forceUpdate = forceUpdate
	conf.platform = platform
	conf.reloadReady = make(chan struct{})

	// we are making sure that warpdrive folder does exist in document path
	os.MkdirAll(filepath.Join(documentPath, "warpdrive"), os.ModePerm)

	// load warpFile
	warpFile, err := getWarpFile()
	if err != nil {
		return err
	}

	conf.warpFile = warpFile

	// we need to create api
	conf.api = makeAPI(warpFile)

	return nil
}

// DownloadVersion it starts download and save the bundle inside path
// it uses the version
func DownloadVersion(cycleID int64, version string) error {
	appID := conf.warpFile.App.ID

	r, err := conf.api.downloadVersion(appID, cycleID, version, conf.platform)
	if err != nil {
		return err
	}

	return extractBundleStream(appID, cycleID, version, r)
}

// DownloadRelease it starts download and save the bundle inside path
// it uses the release id
func DownloadRelease(cycleID, releaseID int64, version string) error {
	appID := conf.warpFile.App.ID

	r, err := conf.api.downloadRelease(appID, cycleID, releaseID)
	if err != nil {
		return err
	}

	return extractBundleStream(appID, cycleID, version, r)
}

// Latest returns an string version of release map
func Latest(cycleID int64) (string, error) {
	releaseMap, err := releaseMap(cycleID)
	if err != nil {
		return "", err
	}

	jsonString, err := json.Marshal(releaseMap)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
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

	// start the process but only for defaultCyle and if forceUpdate is enabled
	if conf.forceUpdate {
		defer func() {
			go func() {
				// we are blocking here until reload is ready to be called
				// this is a trick which make sure the sync between threads
				<-conf.reloadReady

				// we need to check whether warpify configured to be forceUpdated
				if !conf.forceUpdate {
					return
				}

				// we need to find the default cycle config
				defaultCycleConfig, err := conf.warpFile.GetCycle(conf.defaultCycle)
				if err != nil {
					return
				}

				releaseMap, err := releaseMap(defaultCycleConfig.ID)
				if err != nil {
					return
				}

				if softRelease, ok := releaseMap["soft"]; ok && softRelease != nil {
					// start downloading the release
					err = DownloadRelease(defaultCycleConfig.ID, softRelease.ID, softRelease.Version)
					if err != nil {
						return
					}

					// We need to call the native to restart the app
					Reload(defaultCycleConfig.ID, softRelease.Version)
				}
			}()
		}()
	}

	// So we grab the ActiveCycle to find the current active cycle
	// then we look for current version
	// if current version is bundle then, we return "" and let the
	// native code handle the URL
	// if not then we need to generate a path to new value

	version := versionMap.Version(versionMap.ActiveCycle)

	// we need to search for current version inside avaiable map
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
	return filepath.Join(warpBundlePath(appID, cycleConfig.ID, version.Current), "main.jsbundle")
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
	versionMap := &config.VersionMap{}

	if exists, _ := folder.PathExists(config.VersionPath(conf.documentPath)); exists {
		err := versionMap.Load(conf.documentPath)
		if err != nil {
			return nil, err
		}

		// we need to check for one more thing, what if we publish something into app store or
		// play store and user download a new update from there. the versions still be there
		// but it needs to be clean up, this check does this for us, it will clean the warpdrive folder and
		// rebuild everything from scratch if and only if bundle version is differeant than bundle version
		// inside versions.warp file

		if versionMap.BundleVersion(conf.defaultCycle) == conf.bundleVersion {
			return versionMap, nil
		}
	}

	// this is the first time we are creating versions.warp or we need to override it
	// this file will be saved in documentPath/warpdrive/versions.warp
	versionMap = &config.VersionMap{}

	version, err := semver.Make(conf.bundleVersion)
	if err != nil {
		return nil, err
	}

	// if the bundle contains pre value then it consider
	// default cycle name, make sure in your code, defaultCycle and
	// pre in bundle version was the same
	cycleName := cycleNameFromVersion(&version)

	// set the first version inside versions.warp
	// and because the first version points in bundle, then
	// isBundle in `SetCurrentVersion` will be true
	versionMap.ActiveCycle = cycleName
	versionMap.SetCurrentVersion(cycleName, conf.bundleVersion, true)

	// save the file on disk
	err = versionMap.Save(conf.documentPath)
	if err != nil {
		return nil, err
	}

	return versionMap, nil
}

func releaseToString(releases []map[string]*data.Release) string {
	str, _ := json.Marshal(releases)
	return string(str)
}

func warpBundlePath(appID, cycleID int64, version string) string {
	path := fmt.Sprintf("warpdrive/warp.%d.%d.%s", appID, cycleID, version)
	return filepath.Join(conf.documentPath, path)
}

func cycleNameFromVersion(version *semver.Version) string {
	var cycleName string

	if len(version.Pre) > 0 {
		cycleName = version.Pre[0].String()
	} else {
		cycleName = conf.defaultCycle
	}

	return cycleName
}

func extractBundleStream(appID, cycleID int64, version string, r io.ReadCloser) error {
	path := warpBundlePath(appID, cycleID, version)

	// need to make fodler
	err := os.MkdirAll(path, os.ModePerm)
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

func releaseMap(cycleID int64) (map[string]*data.Release, error) {
	appID := conf.warpFile.App.ID

	versionMap, err := getVersionMap()
	if err != nil {
		return nil, err
	}

	cycleConfig, err := conf.warpFile.GetCycleByID(appID, cycleID)
	if err != nil {
		return nil, err
	}

	currentVersion := versionMap.CurrentVersion(cycleConfig.Name)

	releaseMap, err := conf.api.checkVersion(appID, cycleConfig.ID, conf.platform, currentVersion)
	if err != nil {
		return nil, err
	}

	return releaseMap, nil
}
