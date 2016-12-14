package cli

import (
	"fmt"

	"github.com/pressly/warpdrive/config"
	"github.com/spf13/cobra"
)

var (
	listSetup   bool
	resetSetup  bool
	addNewCycle bool
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup warpdrive",
	Long:  `setup warpdrive for react-native project`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isReactNativeProject() {
			return
		}

		globalConf := config.NewGlobalConfig()
		localConf := config.NewClientConfigsForCli()

		if !listSetup && resetSetup {
			localConf.Save()
		}

		err := localConf.Load()
		if err != nil {
			addNewCycle = false
		}

		if listSetup {
			fmt.Println(localConf)
			return
		}

		if !localConf.IsSetupRequired() && !addNewCycle {
			fmt.Println("WarpFile already created for this project, for any modification, use available flags")
			return
		}

		if !addNewCycle || localConf.ServerAddr == "" {
			localConf.ServerAddr = terminalInput("Server address:", false)
		}

		api, err := newAPI(localConf.ServerAddr, globalConf)
		if err != nil {
			fmt.Println(err)
			return
		}

		var appConfig *config.AppConfig

		if !addNewCycle {
			appConfig = &config.AppConfig{}
			appConfig.Name = terminalInput("Enter app name:", false)

			app, err := api.getAppByName(appConfig.Name)
			if err != nil {
				// it means that app does not exists, let's create one
				app, err = api.createApp(appConfig.Name)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}

			// at this point we have app
			appConfig.ID = app.ID
			localConf.App = *appConfig
		} else {
			appConfig = &localConf.App
		}

		for {
			goahead := terminalInputExpect("Add a cycle?", []string{"y", "n"}, "y")
			if goahead == "n" {
				break
			}
			cycleConfig := config.CycleConfig{}
			cycleConfig.Name = terminalInput("Cycle name:", false)

			// try to load the cycle from server
			cycle, err := api.getCycleByName(appConfig.ID, cycleConfig.Name)

			// not found
			if err != nil {
				// we need to create one
				cycle, err = api.createCycle(appConfig.Name, cycleConfig.Name)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
			}

			// at this point we have cycle,
			cycleConfig.ID = cycle.ID
			cycleConfig.Key = cycle.PublicKey

			localConf.AddCycle(&cycleConfig)
		}

		// without cycle, warp is not setup correctly
		if len(localConf.Cycles) == 0 {
			fmt.Println("you need to add at least one cycle")
			return
		}

		// going to save the local, global already saved inside `newAPI` method
		err = localConf.Save()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func initSetupFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&listSetup, "list", "l", false, "list the current setup configuration")
	cmd.Flags().BoolVarP(&resetSetup, "reset", "r", false, "reset configuration of project")
	cmd.Flags().BoolVarP(&addNewCycle, "cycle", "c", false, "add new cycle")
}

func init() {
	initSetupFlags(setupCmd)
	RootCmd.AddCommand(setupCmd)
}
