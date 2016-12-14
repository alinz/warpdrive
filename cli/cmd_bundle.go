package cli

import (
	"fmt"
	"os"

	"path"

	"github.com/pressly/warpdrive/config"
	"github.com/spf13/cobra"
)

// bundle command ...
//////////////////////////

/*
node --max-old-space-size=8192                                                 \
  node_modules/react-native/local-cli/cli.js bundle                            \
  --platform "$PLATFORM"                                                       \
  --entry-file "index.$PLATFORM.js"                                            \
  --bundle-output ./.release/main.jsbundle                                     \
  --assets-dest ./.release                                                     \
  --dev false
*/

var bundlePlatform string

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "bundles react-native project",
	Long:  `bundles react-native project for both ios and android`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isReactNativeProject() {
			return
		}

		// we need this to be really high since big project requires large heap
		iosBundleTask := newCmdTask(
			`ios bundle`,

			`node`,
			`--max-old-space-size=4096`,
			`node_modules/react-native/local-cli/cli.js`,
			`bundle`,
			`--platform`,
			`ios`,
			`--entry-file`,
			`index.ios.js`,
			`--bundle-output`,
			path.Join(config.BundlePathIOS, "main.jsbundle"),
			`--assets-dest`,
			config.BundlePathIOS,
			`--dev`,
			`false`,
		)

		// we need this to be really high since big project requires large heap
		androidBundleTask := newCmdTask(
			`android bundle`,

			`node`,
			`--max-old-space-size=4096`,
			`node_modules/react-native/local-cli/cli.js`,
			`bundle`,
			`--platform`,
			`android`,
			`--entry-file`,
			`index.android.js`,
			`--bundle-output`,
			path.Join(config.BundlePathAndroid, "main.jsbundle"),
			`--assets-dest`,
			config.BundlePathAndroid,
			`--dev`,
			`false`,
		)

		var tasks []task

		switch bundlePlatform {
		case "ios":
			os.RemoveAll(config.BundlePathIOS)
			os.MkdirAll(config.BundlePathIOS, os.ModePerm)
			tasks = append(tasks, iosBundleTask)
		case "android":
			os.RemoveAll(config.BundlePathAndroid)
			os.MkdirAll(config.BundlePathAndroid, os.ModePerm)
			tasks = append(tasks, androidBundleTask)
		case "all":
			os.RemoveAll("./.bundles")
			os.MkdirAll(config.BundlePathIOS, os.ModePerm)
			os.MkdirAll(config.BundlePathAndroid, os.ModePerm)
			tasks = append(tasks, iosBundleTask, androidBundleTask)
		default:
			fmt.Println("platfrom is not recongnized")
			return
		}

		ok, errors := taskRunner(tasks...)

		if !ok {
			for _, err := range errors {
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			return
		}
	},
}

func initFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&bundlePlatform, "platform", "p", "all", "target platform")
}

func init() {
	initFlags(bundleCmd)
	RootCmd.AddCommand(bundleCmd)
}
