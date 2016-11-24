package cli

import "github.com/mitchellh/cli"

type initCommand struct {
	ui cli.Ui
}

func (t *initCommand) Help() string {
	return "setup WarpFile"
}

func (t *initCommand) Run(args []string) int {
	projectConfig, err := ProjectConfig()
	if err == nil {
		return 0
	}

	serverAddr := Input("warpdriver server address:", false)
	projectConfig.setServerAddr(serverAddr)

	err = configSave(projectConfig)
	if err != nil {
		t.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (t *initCommand) Synopsis() string {
	return "Setup WarpFile for the first time in current folder"
}

func newInitCommandFn(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &initCommand{
			ui: ui,
		}, nil
	}
}
