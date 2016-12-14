package cli

import "github.com/spf13/cobra"

// release command ...
//////////////////////////

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "release commands",
	Long:  `create new release, add and remove release from app`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)
}
