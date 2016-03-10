package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/config"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/web"
	"github.com/pressly/warpdrive/web/security"

	"github.com/tylerb/graceful"
)

var (
	flags    = flag.NewFlagSet("warpdrive", flag.ExitOnError)
	confFile = flags.String("config", "", "path to config file")
)

func main() {
	flags.Parse(os.Args[1:])

	//load configuration
	if config, err := config.Load(*confFile, os.Getenv("CONFIG")); err != nil {
		log.Fatal(err)
	} else {
		warpdrive.Config = config
	}

	//initializae database
	if db, err := data.InitDbWithConfig(warpdrive.Config); err != nil {
		log.Fatal(err)
	} else {
		warpdrive.DB = db
	}

	//setup web security such as jwt
	security.SetupWebSecurity()

	//launch the api
	graceful.Run(warpdrive.Config.Server.Bind, 10*time.Second, web.New())
}
