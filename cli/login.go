package cli

import "github.com/mitchellh/cli"
import "fmt"

type loginCommand struct {
	ui cli.Ui
}

func (l *loginCommand) Help() string {
	return "login into warpdrive"
}

func (l *loginCommand) Run(args []string) int {
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

	url, err := apiURL("/session/start")
	if err != nil {
		l.ui.Error(err.Error())
		return 1
	}

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
		l.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (l *loginCommand) Synopsis() string {
	return "comes with no argument, will log you in to warpdrive"
}

func newLoginCommandFn(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &loginCommand{
			ui: ui,
		}, nil
	}
}
