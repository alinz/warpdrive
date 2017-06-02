package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
)

func main() {
	commandEnv := &struct {
		CA   string `require:"true"`
		Crt  string `require:"true"`
		Key  string `require:"true"`
		Port string `require:"true"`
	}{}

	err := envconfig.Process("command", commandEnv)
	if err != nil {
		log.Fatal(err.Error())
	}

	queryEnv := &struct {
		CA   string `require:"true"`
		Crt  string `require:"true"`
		Key  string `require:"true"`
		Port string `require:"true"`
	}{}

	err = envconfig.Process("query", queryEnv)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcCommandConfig, err := helper.NewGrpcConfig(commandEnv.CA, commandEnv.Crt, commandEnv.Key)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcQueryConfig, err := helper.NewGrpcConfig(queryEnv.CA, queryEnv.Crt, queryEnv.Key)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcCommandServer, err := grpcCommandConfig.CreateServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcQueryServer, err := grpcQueryConfig.CreateServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := openDB("/db/warpdrive.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	pb.RegisterCommandServer(grpcCommandServer, &commandServer{db})
	lnCommand, err := net.Listen("tcp", fmt.Sprintf(":%s", commandEnv.Port))
	if err != nil {
		log.Fatal(err.Error())
	}

	pb.RegisterQueryServer(grpcQueryServer, &queryServer{db})
	lnQuery, err := net.Listen("tcp", fmt.Sprintf(":%s", queryEnv.Port))
	if err != nil {
		log.Fatal(err.Error())
	}

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
