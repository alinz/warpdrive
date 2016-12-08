package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
	login logic

	check if we have local config
		- yes:
			- check if local config is new (just initlized)
				- yes:
					- ask for server url
					- display server url
					- ask for email and password
					- try to login
						- success:
							- save the seesion with server address to global config
							- update only server address of local config
							- done
						- failure:
							- show error message that email or password are incorrect
							- done
				- no:
					- load the server url from local config
					- check with global config for existing session assign for that server url
						- yes:
							- validate the session
								- sucess:
									- do nothign
									- done
								- failure:
									- display server url
									- ask for email and password
									- try to login
										- success:
											- save the seesion with server address to global config
											- update only server address of local config
											- done
										- failure:
											- show error message that email or password are incorrect
											- done
						- no:
							- display server url
							- ask for email and password
							- try to login
								- success:
									- save the seesion with server address to global config
									- update only server address of local config
									- done
								- failure:
									- show error message that email or password are incorrect
									- done
		- no:
			- ask for server url
			- display server url
			- ask for email and password
			- try to login
				- success:
					- save the seesion with server address to global config
					- done
				- failure:
					- show error message that email or password are incorrect
					- done
*/

func intractiveLogin(serverAddr string) (string, string, error) {
	if serverAddr == "" {
		serverAddr = terminalInput("Server Address:", false)
	}

	fmt.Printf("Login to %s\n", serverAddr)
	email := terminalInput("Email:", false)
	password := terminalInput("Password:", true)

	api := newAPI(serverAddr)

	err := api.login(email, password)
	if err != nil {
		return "", "", err
	}

	return serverAddr, api.session, nil
}

func validateSession(serverAddr, session string) error {
	api := newAPI(serverAddr)
	api.session = session
	return api.validate()
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login into server",
	Long:  `try to load the local WarpFile for server path. need to type the server path`,
	Run: func(cmd *cobra.Command, args []string) {

		// appID, err := terminalInputAsInt64("App id:", false)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return
		// }

		var session string
		var serverAddr string
		var err error

		// load local config for the current path project
		var globalConfig globalConfig
		var localConfig localConfig

		// we simply load the global config, we don't really care
		// about whether it fails or not, since we are going to write a new one anyway
		globalConfig.Load()

		err = localConfig.Load()
		if err != nil {
			serverAddr = terminalInput("Server Address:", false)
			session, err = globalConfig.getSessionFor(serverAddr)
			if err == nil {
				// we need to check whether the session is valid or not and if it is not,
				// we need to login the user
				err = validateSession(serverAddr, session)
				if err == nil {
					return
				}
			}

			// if the user is not in a project folder, then it means that
			// user just wants to login into a warpdrive server
			serverAddr, session, err = intractiveLogin(serverAddr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			globalConfig.setSessionFor(serverAddr, session)
			globalConfig.Save()
		} else {
			// the local config is new and need to be setup
			if localConfig.isRequiredSetup() {
				serverAddr, session, err = intractiveLogin("")
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				localConfig.ServerAddr = serverAddr
				localConfig.Save()

				globalConfig.setSessionFor(serverAddr, session)
				globalConfig.Save()
			} else {
				// the local config is not new and we need to know
				// whether user has a session for configured serverAddr
				session, err = globalConfig.getSessionFor(localConfig.ServerAddr)
				if err == nil {
					// we need to check whether the session is valid or not and if it is not,
					// we need to login the user
					err = validateSession(localConfig.ServerAddr, session)
				}

				// so we don't have the seesion or session is not valid anymore,
				// so we need to login the user
				if err != nil {
					serverAddr, session, err = intractiveLogin(localConfig.ServerAddr)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					globalConfig.setSessionFor(serverAddr, session)
					globalConfig.Save()
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
