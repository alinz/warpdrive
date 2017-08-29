package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pressly/warpdrive/server"
)

const usagestr = `
Usage: server [options]

Warpdrive Server Options:
    -c, --config <filename>     warpdrive server configuration
`

func usage() {
	fmt.Printf("%s\n", usagestr)
	os.Exit(0)
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "server configuration file")
	flag.Parse()

	if configPath == "" {
		usage()
	}

	server, err := server.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
