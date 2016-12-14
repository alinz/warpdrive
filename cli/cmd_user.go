package cli

import "github.com/spf13/cobra"

// user command ...
//////////////////////////

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "user commands",
	Long:  `create new user, add and remove user from app`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
