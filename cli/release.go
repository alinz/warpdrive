package cli

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

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

// grab all the files from given path, included nested folder as well
func allFilesForPath(path string) ([]string, error) {
	var files []string

	// we need to define this loop to make it available inside loop itself
	var loop func(string, *[]string) error

	loop = func(path string, files *[]string) error {
		dir, err := os.Open(path)
		if err != nil {
			return err
		}

		defer dir.Close()

		dirStat, err := dir.Stat()
		if err != nil {
			return err
		}

		if !dirStat.IsDir() {
			*files = append(*files, path)
			return nil
		}

		fileInfos, err := dir.Readdir(-1)
		if err != nil {
			return err
		}

		for _, fileInfo := range fileInfos {
			err = loop(filepath.Join(path, fileInfo.Name()), files)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := loop(path, &files)

	return files, err
}

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
		return nil, err
	}

	bundleFiles, err := allFilesForPath(path)
	if err != nil {
		return nil, err
	}

	bundleFilesMap := make(map[string]string)
	for _, file := range bundleFiles {
		bundleFilesMap[file] = file
	}

	r, w := io.Pipe()

	go func() {
		writer := multipart.NewWriter(w)
		defer writer.Close()

		partWriter, err := writer.CreateFormFile("file", "bundle.tar.gz")
		if err != nil {
			return
		}

		err = warp.Compress(bundleFilesMap, partWriter)

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
		var r io.Reader

		// creating an api to make sure the information about
		// account in current react-native project is correct
		// before doing long process of bundling.
		api, err := getActiveAPI(nil, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		switch publishPlatform {
		case "ios", "android":
			r, err = bundleReader(publishPlatform)
		case "all":
			r, err = bundleReader("ios")
			if err != nil {
				break
			}
			r, err = bundleReader("android")
		default:
			err = errPlatformBundleNotRecognized
		}

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		bundles, err := api.bundleUpload(0, 0, 0, r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(bundles)
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