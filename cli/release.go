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

		// check if localConfig is new
		if !local.isRequiredSetup() {
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

// app ...

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
	initReleaseConfigureFlags(releaseConfigureCmd)
	releaseCmd.AddCommand(releaseConfigureCmd)

	// adding release command to Root
	RootCmd.AddCommand(releaseCmd)
}
