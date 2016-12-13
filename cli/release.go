package cli

import (
	"fmt"
	"io"
	"strings"

	"path/filepath"

	"github.com/blang/semver"
	"github.com/pressly/warpdrive/constants"
	"github.com/pressly/warpdrive/lib/folder"
	"github.com/pressly/warpdrive/lib/warp"
	"github.com/spf13/cobra"
)

// release configure ...
//////////////////////////

var releaseConfigureSwitch bool

var releaseConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configures the project",
	Long: `
configures the current project with existing app and cycle.
this needs to be called once per project, unless the app or cycle or both have been removed from warpdrive server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if we are in react-native project
		if !isReactNativeProject() {
			fmt.Println("the current path is not a react-native project")
			return
		}

		var err error

		var global globalConfig
		var local localConfig

		global.Load()
		err = local.Load()
		if err != nil {
			fmt.Println("please call 'warp init' before calling this command")
			return
		}

		// check if localConfig is new and it will override and skip it if
		// releaseConfigureSwitch is being set
		if !local.isRequiredSetup() && !releaseConfigureSwitch {
			fmt.Println("the current project has already been setup, if you want to change, use --switch flag")
			return
		}

		api, err := getActiveAPI(&global, &local)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		appName := terminalInput("App's name:", false)
		cycleName := terminalInput("Cycle's name:", false)

		// making sure user can access app
		app, err := api.getAppByName(appName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// make sure the cycleName exists
		cycle, err := api.getCycleByName(app.ID, cycleName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// update the local config file and save it back
		local.AppID = app.ID
		local.CycleID = cycle.ID
		local.Key = cycle.PublicKey
		local.ServerAddr = api.serverAddr

		local.Save()
	},
}

func initReleaseConfigureFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&releaseConfigureSwitch, "switch", "s", false, "request configure switch")
}

// release publish ...
//////////////////////////

var publishPlatform string
var publishVersion string
var publishNote string

var errIosBundleNotFound = fmt.Errorf("ios bundle not found")
var errAndroidBundleNotFound = fmt.Errorf("android bundle not found")
var errPlatformBundleNotRecognized = fmt.Errorf("platform not recognized")
var errPathIsNotDir = fmt.Errorf("path is not a directory")
var errVersionBundleNotSet = fmt.Errorf("version need to be set for bundle")
var errVersionBundleFormatInvalid = fmt.Errorf("verssion bundle format is invalid")

func isBundleReady(platform string) bool {
	return true
}

// bundleReader will tar and gzip the entire bundle as stream
// and returns an io.Reader. This is an optimization to
// reduce the memory usgae during sending the large bundle
func bundleReader(platform string) (io.Reader, error) {
	var err error
	var path string

	// setup the path for given platform
	switch platform {
	case "ios":
		if !isBundleReady("ios") {
			err = errIosBundleNotFound
		} else {
			path = constants.BundlePathIOS
		}
	case "android":
		if !isBundleReady("android") {
			err = errAndroidBundleNotFound
		} else {
			path = constants.BundlePathAndroid
		}
	default:
		err = errPlatformBundleNotRecognized
	}

	if err != nil {
		return nil, err
	}

	bundleFiles, err := folder.ListFilePaths(path)
	if err != nil {
		return nil, err
	}

	cleanedPath := filepath.Clean(path) + "/"

	bundleFilesMap := make(map[string]string)
	for _, file := range bundleFiles {
		bundleFilesMap[strings.Replace(file, cleanedPath, "", 1)] = file
	}

	r, w := io.Pipe()

	go func() {
		err = warp.Compress(bundleFilesMap, w)
		// we need to close the pipe's writer to
		// signal the reader that there won't be any more bytes coming
		if err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()

	return r, nil
}

// publishPlatform must be ios, android
// publishVersion must be auto or x.y.z format
// publishNote is optional and describe the release note for this version
func prepareRelease(api *api, appID, cycleID int64, platform, version, note string) (io.Reader, int64, error) {
	reader, err := bundleReader(platform)
	if err != nil {
		return nil, 0, err
	}

	release, err := api.createRelease(appID, cycleID, platform, version, note)
	if err != nil {
		return nil, 0, err
	}

	return reader, release.ID, nil
}

var releasePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish the project",
	Long: `
publish the current bundle projects, ios and android, to warpdrive server
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isReactNativeProject() {
			fmt.Println("current path is not react-native project")
			return
		}

		// check if version passed as an argument
		if publishVersion == "" {
			fmt.Println(errVersionBundleNotSet.Error())
			return
		}

		// just make sure that version is corectly formated before sending it to server
		_, err := semver.Make(publishVersion)
		if err != nil {
			fmt.Println(errVersionBundleFormatInvalid.Error())
			return
		}

		var iosBundle io.Reader
		var androidBundle io.Reader
		var iosReleaseID int64
		var androidReleaseID int64

		var global globalConfig
		var local localConfig

		// creating an api to make sure the information about
		// account in current react-native project is correct
		// before doing long process of bundling.
		api, err := getActiveAPI(&global, &local)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// if we get here then, we can use global and local varibales
		// those are fully loaded when being passed to `getActiveAPI`
		// so now, we need to get the `app.id`, `cycle.id` from local config
		appID := local.AppID
		cycleID := local.CycleID

		switch publishPlatform {
		case "ios":
			iosBundle, iosReleaseID, err = prepareRelease(api, appID, cycleID, "ios", publishVersion, publishNote)
		case "android":
			androidBundle, androidReleaseID, err = prepareRelease(api, appID, cycleID, "android", publishVersion, publishNote)
		case "all":
			iosBundle, iosReleaseID, err = prepareRelease(api, appID, cycleID, "ios", publishVersion, publishNote)
			if err != nil {
				break
			}
			androidBundle, androidReleaseID, err = prepareRelease(api, appID, cycleID, "android", publishVersion, publishNote)
		default:
			err = errPlatformBundleNotRecognized
		}

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if iosBundle != nil {
			_, err = api.bundleUpload(appID, cycleID, iosReleaseID, iosBundle)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("upload ios bundle is completed")
		}

		if androidBundle != nil {
			_, err = api.bundleUpload(appID, cycleID, androidReleaseID, androidBundle)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("upload android bundle is completed")
		}

		// we just need to lock the release for any releases with correct id
		if iosReleaseID != 0 {

		}

		if androidReleaseID != 0 {

		}
	},
}

func initReleasePublishFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&publishPlatform, "platform", "p", "all", "publish specific platform, [ios, android]")
	cmd.Flags().StringVarP(&publishVersion, "version", "v", "", "publish version. use semantic versioning x.y.z")
	cmd.Flags().StringVarP(&publishNote, "note", "n", "", "add release note to new version")
}

// release ...
//////////////////////////

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "configure and publish new releases",
	Long:  `configure and publish new releases to warpdrive server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("please run 'warp release -h' for help")
	},
}

func initReleaseFlags(cmd *cobra.Command) {

}

func init() {
	initReleaseFlags(releaseCmd)

	// commands under release
	// configure
	initReleaseConfigureFlags(releaseConfigureCmd)
	releaseCmd.AddCommand(releaseConfigureCmd)
	// publish
	initReleasePublishFlags(releasePublishCmd)
	releaseCmd.AddCommand(releasePublishCmd)

	// adding release command to Root
	RootCmd.AddCommand(releaseCmd)
}
