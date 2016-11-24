package cli

import (
	"strings"

	"github.com/mitchellh/cli"
)

type tempCommand struct {
	ui cli.Ui
}

func (t *tempCommand) Help() string {
	helpText := `
	
	
	`

	return strings.TrimSpace(helpText)
}

func (t *tempCommand) Run(args []string) int {
	return 0
}

func (t *tempCommand) Synopsis() string {
	return "one line help"
}

func newTempCommandFn(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &tempCommand{
			ui: ui,
		}, nil
	}
}
