package cli

import "github.com/mitchellh/cli"
import "fmt"

type loginCommand struct{}

func (t *loginCommand) Help() string {
	return "login into warpdrive"
}

func (t *loginCommand) Run(args []string) int {
	projectConfig, err := ProjectConfig()
	if err != nil {
		fmt.Println("WarpFile not found")
		return 1
	}

	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    Input("email:", false),
		Password: Input("password", true),
	}

	url := apiURL("/session/start")
	resp, err := httpRequest("POST", url, reqBody)

	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		return 1
	}

	cookie := resp.Cookies()[0]
	jwt := cookie.Value

	globalConfig, _ := GlobalConfig()
	globalConfig.setSession(projectConfig.Server.Addr, jwt)

	err = configSave(globalConfig)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	return 0
}

func (t *loginCommand) Synopsis() string {
	return "comes with no argument, will log you in to warpdrive"
}

func LoginCommand() cli.CommandFactory {
	return func() (cli.Command, error) {
		return &loginCommand{}, nil
	}
}
