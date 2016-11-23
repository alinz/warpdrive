package main

import (
	"log"
	"os"

	cliLib "github.com/mitchellh/cli"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/cli"
)

func main() {
	c := cliLib.NewCLI("warp", warpdrive.VERSION)

	c.Args = os.Args[1:]
	c.Commands = cli.Commands

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
