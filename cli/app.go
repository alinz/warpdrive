package cli

import (
	"fmt"

	"github.com/pressly/warpdrive"
	"github.com/spf13/cobra"
)

// app create ...

var appName string

var appCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a new app",
	Long:  `creates a new app`,
	Run: func(cmd *cobra.Command, args []string) {
		if appName == "" {
			fmt.Println("please provide the app's name")
			return
		}

		api, err := getActiveAPI()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err = api.createApp(appName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func initAppCreateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&appName, "name", "n", "", "app's name")
}

// app ...

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "working with app",
	Long:  `working with apps in worpdrive`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Warp %s\n", warpdrive.VERSION)
	},
}

func initAppFlags(cmd *cobra.Command) {

}

func init() {
	initAppFlags(appCmd)

	// commands under app
	initAppCreateFlags(appCreateCmd)
	appCmd.AddCommand(appCreateCmd)

	// adding app command to Root
	RootCmd.AddCommand(appCmd)
}
