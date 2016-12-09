package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// release configure ...

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

// release publish

var releasePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish the project",
	Long: `
publish the current bundle projects, ios and android, to warpdrive server
`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func initReleasePublishFlags(cmd *cobra.Command) {
}

// release ...

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
