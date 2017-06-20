package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/asdine/storm"
	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
	"github.com/pressly/warpdrive/server"
	"github.com/pressly/warpdrive/server/config"
)

const usagestr = `
Usage: server [options]

server Options:
    -c, --config <filename>     warpdrive server's configuration
`

func usage() {
	fmt.Printf("%s\n", usagestr)
	os.Exit(0)
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "service configuration file")
	flag.Parse()

	conf := &config.Config{}
	err := helper.ConfigFile(configPath, conf)
	if err != nil {
		log.Fatal(err)
	}

	// open database connection
	db, err := storm.Open(conf.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	// create grpc server for Command Service
	cmdConf := conf.Command
	grpcCommandConfig, err := helper.NewGrpcConfig(cmdConf.CAPath, cmdConf.CertPath, cmdConf.KeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcCommandServer, err := grpcCommandConfig.CreateServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	pb.RegisterCommandServer(grpcCommandServer, server.NewCommandServer(db, conf))
	lnCommand, err := net.Listen("tcp", cmdConf.Addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	// create grpc for Query Service
	qryConf := conf.Query
	grpcQueryConfig, err := helper.NewGrpcConfig(qryConf.CAPath, qryConf.CertPath, qryConf.KeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcQueryServer, err := grpcQueryConfig.CreateServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	pb.RegisterQueryServer(grpcQueryServer, server.NewQueryServer(db, conf))
	lnQuery, err := net.Listen("tcp", qryConf.Addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	// run both Command and Query services in their own goroutines
	commandCloseChan := make(chan error)
	go func() {
		commandCloseChan <- grpcCommandServer.Serve(lnCommand)
	}()

	queryCloseChan := make(chan error)
	go func() {
		queryCloseChan <- grpcQueryServer.Serve(lnQuery)
	}()

	// proper graceful shutdown of services
	// this select waits until one of the services
	// sends a nil or error. In either cases, we need to
	// shutdown the other service gracefully and log the error
	select {
	case err := <-commandCloseChan:
		if err != nil {
			log.Print(err.Error())
		}
		grpcQueryServer.GracefulStop()
	case err := <-queryCloseChan:
		if err != nil {
			log.Print(err.Error())
		}
		grpcCommandServer.GracefulStop()
	}
}
