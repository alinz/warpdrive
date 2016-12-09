package cli

import (
	"fmt"
	"os"

	"path/filepath"

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

// grab all the files from given path, included nested folder as well
func allFilesForPath(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			files = append(files, path)
		}
		return err
	})

	return files, err
}

func isBundleReady(platform string) bool {
	return true
}

func uploadBundleFor(platform string) error {
	var err error

	var path string

	// setup the path for given platform
	switch platform {
	case "ios":
		if !isBundleReady("ios") {
			err = errIosBundleNotFound
		} else {
			path = iosBundlePath
		}
	case "android":
		if !isBundleReady("android") {
			err = errAndroidBundleNotFound
		} else {
			path = androidBundlePath
		}
	default:
		err = errPlatformBundleNotRecognized
	}

	if err != nil {
		return err
	}

	bundleFiles, err := allFilesForPath(path)
	if err != nil {
		return err
	}

	fmt.Println(bundleFiles)

	return nil
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

		var err error

		switch publishPlatform {
		case "ios", "android":
			err = uploadBundleFor(publishPlatform)
		case "all":
			err = uploadBundleFor("ios")
			if err != nil {
				break
			}
			err = uploadBundleFor("android")
		default:
			err = errPlatformBundleNotRecognized
		}

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func initReleasePublishFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&publishPlatform, "platform", "p", "all", "publish specific platform, [ios, android]")
	cmd.Flags().StringVarP(&publishVersion, "version", "v", "auto", "publish version. use semantic versioning x.y.z")
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
