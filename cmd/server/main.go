package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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

	start, err := server.SetupServer()
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", server.Conf.Server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	err = start(ln)
	if err != nil {
		log.Fatal(err)
	}
}
