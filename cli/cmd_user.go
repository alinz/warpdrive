package cli

import "github.com/spf13/cobra"

// user command ...
//////////////////////////

// warp user -a app -u ali@pressly.com
// warp user add -a app
// warp user rm -a app -u ali@pressly.com
// warp user create -e "" -p ""

var userEmailFlag string
var userPasswordFlag string
var userAppNameFlag string

/////////////// Create User Command

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new user",
	Long:  `create a new user`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initUserCreateFlags(cmd *cobra.Command) {

}

/////////////// Add User Command

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a user to app",
	Long:  `add a user to app, so they can release and configure app`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initUserAddFlags(cmd *cobra.Command) {

}

/////////////// Remove User Command

var userRemoveCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove a user form app",
	Long:  `remove a user from app, permission wise, user no longer can read or publish new releases`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initUserRemoveFlags(cmd *cobra.Command) {

}

/////////////// Root User Command

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "user commands",
	Long:  `create new user, add and remove user from app`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initUserFlags(cmd *cobra.Command) {

}

func init() {
	initUserCreateFlags(userCreateCmd)
	userCmd.AddCommand(userCreateCmd)

	initUserAddFlags(userAddCmd)
	userCmd.AddCommand(userAddCmd)

	initUserRemoveFlags(userRemoveCmd)
	userCmd.AddCommand(userRemoveCmd)

	initUserFlags(userCmd)
	RootCmd.AddCommand(userCmd)
}
