package cli

import "github.com/mitchellh/cli"

var Commands map[string]cli.CommandFactory

func init() {
	Commands = make(map[string]cli.CommandFactory)

	Commands["init"] = InitCommand()
	Commands["login"] = LoginCommand()
}
