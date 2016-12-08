package cli

import (
	"fmt"

	"github.com/pressly/warpdrive"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of warp",
	Long:  `display the current version of warpdrive cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Warp %s\n", warpdrive.VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
