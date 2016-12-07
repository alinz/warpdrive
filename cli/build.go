package cli

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

/*
node --max-old-space-size=8192                                                 \
  node_modules/react-native/local-cli/cli.js bundle                            \
  --platform "$PLATFORM"                                                       \
  --entry-file "index.$PLATFORM.js"                                            \
  --bundle-output ./.release/main.jsbundle                                     \
  --assets-dest ./.release                                                     \
  --dev false
*/

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "bundles react-native project",
	Long:  `bundles react-native project for both ios and android`,
	Run: func(cmd *cobra.Command, args []string) {
		if isReactNativeProject() {
			fmt.Println("The current path is not a react-native project")
			return
		}

		// need to make couple of folders for bundles for both ios and android
		// we need two sets of folders, one for ios and one for android since
		// we can parallel the bundle process
		os.RemoveAll("./.bundles")
		os.MkdirAll("./.bundles/ios", os.ModePerm)
		os.MkdirAll("./.bundles/android", os.ModePerm)

		ok, errors := taskRunner(
			newCmdTask(
				`ios bundle`,

				`node`,
				`--max-old-space-size=4096`, // we need this to be really high since big project requires large heap
				`node_modules/react-native/local-cli/cli.js`,
				`bundle`,
				`--platform`,
				`ios`,
				`--entry-file`,
				`index.ios.js`,
				`--bundle-output`,
				`./.bundles/ios/main.jsbundle`,
				`--assets-dest`,
				`./.bundles/ios`,
				`--dev`,
				`false`,
			),
			newCmdTask(
				`android bundle`,

				`node`,
				`--max-old-space-size=4096`, // we need this to be really high since big project requires large heap
				`node_modules/react-native/local-cli/cli.js`,
				`bundle`,
				`--platform`,
				`android`,
				`--entry-file`,
				`index.android.js`,
				`--bundle-output`,
				`./.bundles/android/main.jsbundle`,
				`--assets-dest`,
				`./.bundles/android`,
				`--dev`,
				`false`,
			),
		)

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

func init() {
	RootCmd.AddCommand(buildCmd)
}
