package cli

import (
	"os"

	"github.com/mitchellh/cli"
)

// Commands list of all registered and available commands
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = make(map[string]cli.CommandFactory)

	Commands["init"] = newInitCommandFn(ui)
	Commands["login"] = newLoginCommandFn(ui)
	//Commands["app"] = newAppCommandFn(ui)
}
