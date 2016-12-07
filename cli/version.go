package cli

import (
	"fmt"

	"github.com/pressly/warpdrive"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of warp",
	Long:  `All software has versions. This is Warp's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Warp v%s\n", warpdrive.VERSION)
	},
}
