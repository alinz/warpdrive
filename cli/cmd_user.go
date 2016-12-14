package cli

import "github.com/spf13/cobra"
import "fmt"
import "github.com/pressly/warpdrive/config"

// user command ...
//////////////////////////

// warp user add -a app -e ali@pressly.com
// warp user rm -a app -e ali@pressly.com
// warp user create -e ali@pressly.com -p 12345

var userNameFlag string
var userEmailFlag string
var userPasswordFlag string
var userAppNameFlag string

/////////////// Create User Command

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new user",
	Long:  `create a new user`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isReactNativeProject() {
			return
		}

		if userNameFlag == "" {
			fmt.Println("name not defined")
			return
		}

		if userEmailFlag == "" {
			fmt.Println("email not defined")
			return
		}

		if userPasswordFlag == "" {
			fmt.Println("password not defined")
			return
		}

		global := config.NewGlobalConfig()
		global.Load()

		local := config.NewClientConfigsForCli()
		err := local.Load()
		if err != nil {
			fmt.Println("run warp setup first")
			return
		}

		api, err := newAPI(local.ServerAddr, global)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err = api.createUser(userNameFlag, userEmailFlag, userPasswordFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func initUserCreateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&userNameFlag, "name", "n", "", "new user's name")
	cmd.Flags().StringVarP(&userEmailFlag, "email", "e", "", "new user's email address")
	cmd.Flags().StringVarP(&userPasswordFlag, "password", "p", "", "news user's password")
}

/////////////// Add User Command

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a user to app",
	Long:  `add a user to app, so they can release and configure app`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isReactNativeProject() {
			return
		}

		if userEmailFlag == "" {
			fmt.Println("email not defined")
			return
		}

		if userAppNameFlag == "" {
			fmt.Println("app not defined")
			return
		}

		global := config.NewGlobalConfig()
		global.Load()

		local := config.NewClientConfigsForCli()
		err := local.Load()
		if err != nil {
			fmt.Println("run warp setup first")
			return
		}

		api, err := newAPI(local.ServerAddr, global)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		user, err := api.getUserByEmail(userEmailFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		app, err := api.getAppByName(userAppNameFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = api.addUserToApp(user.ID, app.ID)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func initUserAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&userEmailFlag, "name", "n", "", "user's email")
	cmd.Flags().StringVarP(&userAppNameFlag, "email", "e", "", "app's name")
}

/////////////// Remove User Command

var userRemoveCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove a user form app",
	Long:  `remove a user from app, permission wise, user no longer can read or publish new releases`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isReactNativeProject() {
			return
		}

		if userEmailFlag == "" {
			fmt.Println("email not defined")
			return
		}

		if userAppNameFlag == "" {
			fmt.Println("app not defined")
			return
		}

		global := config.NewGlobalConfig()
		global.Load()

		local := config.NewClientConfigsForCli()
		err := local.Load()
		if err != nil {
			fmt.Println("run warp setup first")
			return
		}

		api, err := newAPI(local.ServerAddr, global)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		user, err := api.getUserByEmail(userEmailFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		app, err := api.getAppByName(userAppNameFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = api.removeUserFromApp(user.ID, app.ID)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func initUserRemoveFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&userEmailFlag, "name", "n", "", "user's email")
	cmd.Flags().StringVarP(&userAppNameFlag, "email", "e", "", "app's name")
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
