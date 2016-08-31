package main

import (
	"flag"
	"os"
	"syscall"
	"time"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/web/routes"

	"github.com/zenazn/goji/graceful"
)

func main() {
	//warpdrive.Logger.Printf("Version: %s", warpdrive.VERSION)
	flags := flag.NewFlagSet("warpdrive", flag.ExitOnError)
	confFile := flags.String("config", "", "path to config file")
	flags.Parse(os.Args[1:])

	//setup config
	//
	conf, err := warpdrive.NewConfig(*confFile)
	if err != nil {
		warpdrive.Logger.Fatal(err)
	}

	//setup database
	//
	_, err = data.NewDatabase()
	if err != nil {
		warpdrive.Logger.Fatal(err)
	}

	//setup routes
	//
	r := routes.New()

	//graceful shutdown
	//
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	graceful.Timeout(10 * time.Second) // Wait timeout for handlers to finish.
	graceful.PreHook(func() {
		warpdrive.Logger.Println("waiting for requests to finish..")
	})
	graceful.PostHook(func() {
		warpdrive.Logger.Println("...")
	})

	warpdrive.Logger.Printf("Warodrive API server runs at %s\n", conf.Server.Addr)
	err = graceful.ListenAndServe(conf.Server.Addr, r)
	if err != nil {
		warpdrive.Logger.Fatal(err)
	}

	graceful.Wait()
}
