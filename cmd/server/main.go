package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/web/routes"
	"github.com/zenazn/goji/graceful"
)

func main() {
	flags := flag.NewFlagSet("warpdrive", flag.ExitOnError)
	confFile := flags.String("config", "", "path to config file")
	flags.Parse(os.Args[1:])

	// Config
	conf, err := warpdrive.NewConfig(*confFile)
	if err != nil {
		log.Fatal(err)
	}

	// Database
	_, err = data.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	r := routes.New()

	// Server boot
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	graceful.Timeout(10 * time.Second) // Wait timeout for handlers to finish.
	graceful.PreHook(func() {
		log.Println("waiting for requests to finish..")
	})
	graceful.PostHook(func() {
		log.Println("...")
	})

	log.Printf("Warpdrive API server runs at %s\n", conf.Server.Addr)
	err = graceful.ListenAndServe(conf.Server.Addr, r)
	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}
