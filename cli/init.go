package cli

import "github.com/mitchellh/cli"

type initCommand struct{}

func (t *initCommand) Help() string {
	return "setup WarpFile"
}

func (t *initCommand) Run(args []string) int {
	projectConfig, err := ProjectConfig()
	if err == nil {
		return 0
	}

	projectConfig.Server.Addr = Input("warpdriver server address:", false)
	err = configSave(projectConfig)
	if err != nil {
		return 1
	}

	return 0
}

func (t *initCommand) Synopsis() string {
	return "Setup WarpFile for the first time in current folder"
}

func InitCommand() cli.CommandFactory {
	return func() (cli.Command, error) {
		return &initCommand{}, nil
	}
}
