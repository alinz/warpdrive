package cli

import (
	"fmt"

	"github.com/pressly/warpdrive"
	"github.com/spf13/cobra"
)

var appName string

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "working with app",
	Long:  `working with apps in worpdrive`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Warp %s\n", warpdrive.VERSION)
	},
}

func initAppFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&appName, "name", "n", "", "app's name")
}

func init() {
	initAppFlags(appCmd)
	RootCmd.AddCommand(appCmd)
}
