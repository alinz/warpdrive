package cli

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/pressly/warpdrive/config"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/folder"
	"github.com/pressly/warpdrive/lib/warp"
	"github.com/spf13/cobra"
)

// publish command ...
//////////////////////////

// warp publish -v 1.2.3-dev+build12
// warp publish -v 1.2.3-dev+build12 -p ios
// warp publish -v 1.2.3-dev+build12 -p ios -n "first release"

var publishPlatformFlag string
var publishVersionFlag string
var publishNoteFlag string

// bundleReader will tar and gzip the entire bundle as stream
// and returns an io.Reader. This is an optimization to
// reduce the memory usgae during sending the large bundle
func bundleReader(platform string) (io.Reader, error) {
	// setup the path for given platform
	path, err := bundleReadyPath(platform)
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
		bundleFilesMap[file] = strings.Replace(file, cleanedPath, "", 1)
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

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish commands",
	Long:  `upload a new release`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var r io.Reader
		var readers []io.Reader
		var release *data.Release
		var releases []*data.Release

		// make sure the current path is react-native project
		if !isReactNativeProject() {
			return
		}

		// loads global and local configs
		global := config.NewGlobalConfig()
		local := config.NewClientConfigsForCli()

		// we don't really care about global config error
		// at the moment, I use the global to get access to cached session key
		global.Load()

		err = local.Load()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// create an api based on server addr and possible cached session key
		// if not, newApi will ask the user for authentication
		api, err := newAPI(local.ServerAddr, global)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// parse the version, we need to do this to extract cycle name in version
		// version must have prerelease in warpdrive
		version, err := semver.Make(publishVersionFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(version.Pre) == 0 {
			fmt.Println("prerelease not found in version")
			return
		}

		// the first part of prerelease is consider as cycleName
		// 1.2.3-dev.123 => `dev` will be cycleName
		cycleName := version.Pre[0].String()

		// we need to see if local config, configures to access that cycle
		cycleConfig, err := local.GetCycle(cycleName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// we need to clean up the newly created release record if there is an error occuring
		// or we need to lock the release if everything went well
		defer func() {
			if err != nil {
				// if any error happens, we delete all releases that we just added
				for _, release = range releases {
					api.removeRelease(local.App.ID, cycleConfig.ID, release.ID)
				}
			} else {
				// we need to loop through the releases and
				// lock all of them
				for _, release = range releases {
					err = api.lockRelease(local.App.ID, cycleConfig.ID, release.ID)

					// in case of error, we go a little bit aggresive and remove
					// everything that we just added to prevent any corruption
					if err != nil {
						fmt.Println(err.Error())
						// if there is an error, we need to cleaned up the code
						for _, release = range releases {
							api.removeRelease(local.App.ID, cycleConfig.ID, release.ID)
						}
					}
				}
			}
		}()

		switch publishPlatformFlag {
		case "all":
			// creates an ios release record
			release, err = api.createRelease(local.App.ID, cycleConfig.ID, "ios", publishVersionFlag, publishNoteFlag)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			releases = append(releases, release)

			// creates an android release record
			release, err = api.createRelease(local.App.ID, cycleConfig.ID, "android", publishVersionFlag, publishNoteFlag)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			releases = append(releases, release)

			r, err = bundleReader("ios")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			readers = append(readers, r)

			r, err = bundleReader("android")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			readers = append(readers, r)

		default:
			// in default mode, we simply pass the publishPlatformFlag to createRelease to creates a
			// release record. we then append that one to releases slice. That slice will be consumed
			// at the end of the code base, also it will be
			release, err = api.createRelease(local.App.ID, cycleConfig.ID, publishPlatformFlag, publishVersionFlag, publishNoteFlag)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			releases = append(releases, release)

			r, err = bundleReader(publishPlatformFlag)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			readers = append(readers, r)
		}

		for idx, reader := range readers {
			release = releases[idx]
			_, err = api.bundleUpload(local.App.ID, cycleConfig.ID, release.ID, reader)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Printf("published new version for %s", release.Platform.String())
		}
	},
}

func initPublishFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&publishPlatformFlag, "platform", "p", "all", "publish's platform")
	cmd.Flags().StringVarP(&publishVersionFlag, "version", "v", "", "publish's version")
	cmd.Flags().StringVarP(&publishNoteFlag, "note", "n", "", "publish's note")
}

func init() {
	initPublishFlags(publishCmd)
	RootCmd.AddCommand(publishCmd)
}
