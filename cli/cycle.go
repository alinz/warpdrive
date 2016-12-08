package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cycle create ...

var cycleName string

var cycleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a new cycle",
	Long:  `creates a new cycle for app`,
	Run: func(cmd *cobra.Command, args []string) {
		if appName == "" {
			fmt.Println("please provide the app's name")
			return
		}

		if cycleName == "" {
			fmt.Println("please provide the cycle's name")
			return
		}

		api, err := getActiveAPI(nil, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err = api.createCycle(appName, cycleName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func initCycleCreateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&appName, "app", "a", "", "existing app's name")
	cmd.Flags().StringVarP(&cycleName, "name", "n", "", "new cycle's name")
}

// cycle ...

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "working with cycle",
	Long:  `working with cycles in worpdrive`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("please run 'warp cycle create -h' for help")
	},
}

func initCycleFlags(cmd *cobra.Command) {

}

func init() {
	initCycleFlags(cycleCmd)

	// commands under app
	initCycleCreateFlags(cycleCreateCmd)
	cycleCmd.AddCommand(cycleCreateCmd)

	// adding app command to Root
	RootCmd.AddCommand(cycleCmd)
}
