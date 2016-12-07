package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login into server",
	Long:  `try to load the local WarpFile for server path. need to type the server path`,
	Run: func(cmd *cobra.Command, args []string) {

		// load local config for the current path project
		var globalConfig globalConfig
		var localConfig localConfig
		var api *api

		// we simply load the global config, we don't really care
		// about whether it fails or not, since we are going to write a new one anyway
		globalConfig.Load()

		// if the user is not in a project folder, then it means that
		// user just wants to login into a warpdrive server
		err := localConfig.Load()
		if err != nil {
			fmt.Println(err.Error())
			serverAddr := terminalInput("Server Address:", false)
			email := terminalInput("Email:", false)
			password := terminalInput("Password:", true)

			api = newAPI(serverAddr)

			err = api.login(email, password)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			globalConfig.setSessionFor(serverAddr, api.session)
			globalConfig.Save()

			return
		}

		// check if local config is new
		if localConfig.isRequiredSetup() {
			serverAddr := terminalInput("Server Address:", false)

			// we need to check if we have a session in globalConfig
			session, err := globalConfig.getSessionFor(serverAddr)
			if err != nil {
				// it means that this is a new server and we need to ask for email and password
				email := terminalInput("Email:", false)
				password := terminalInput("Password:", true)

				api = newAPI(serverAddr)

				err = api.login(email, password)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				session = api.session
				globalConfig.setSessionFor(serverAddr, api.session)
				globalConfig.Save()
			}

			api.session = session

			// now we need to ask for app id and cycle id and then make sure
			// the current logged in user can access that app
			appID, err := terminalInputAsInt64("App id:", false)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			cycleID, err := terminalInputAsInt64("Cycle id:", false)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			cycle, err := api.getCycle(appID, cycleID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			localConfig.AppID = appID
			localConfig.CycleID = cycleID
			localConfig.Key = cycle.PublicKey

			localConfig.Save()
		} else {
			// we need to see if user has session in globalConfig
			api = newAPI(localConfig.ServerAddr)
			session, err := globalConfig.getSessionFor(localConfig.ServerAddr)
			if err != nil {
				// no session we need to ask the user to login and save the session in globalConfig
				email := terminalInput("Email:", false)
				password := terminalInput("Password:", true)

				err = api.login(email, password)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				session = api.session
				globalConfig.setSessionFor(localConfig.ServerAddr, api.session)
				globalConfig.Save()
			}

			api.session = session

			cycle, err := api.getCycle(localConfig.AppID, localConfig.CycleID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// make sure the public key is the same as server side
			localConfig.Key = cycle.PublicKey

			localConfig.Save()
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
