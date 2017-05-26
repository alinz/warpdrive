package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	context "golang.org/x/net/context"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/kelseyhightower/envconfig"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
)

type grpcConfig struct {
	certPool    *x509.CertPool
	certificate tls.Certificate
}

func (g *grpcConfig) createServer() (*grpc.Server, error) {
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{g.certificate},
		ClientCAs:    g.certPool,
	}

	serverOption := grpc.Creds(credentials.NewTLS(tlsConfig))
	server := grpc.NewServer(serverOption)

	return server, nil
}

func newGrpcConfig(ca, crt, key string) (*grpcConfig, error) {
	certificate, err := tls.LoadX509KeyPair(
		crt,
		key,
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	return &grpcConfig{certPool, certificate}, nil
}

func openDB(path string) (*storm.DB, error) {
	db, err := storm.Open(path, storm.Codec(protobuf.Codec))
	if err != nil {
		return nil, err
	}

	return db, nil
}

type commandServer struct {
	db *storm.DB
}

func (c *commandServer) CreateApp(ctx context.Context, app *pb.App) (*pb.App, error) {
	err := c.db.Save(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (c *commandServer) GetAllApps(_ *pb.Empty, stream pb.Command_GetAllAppsServer) error {
	var apps []pb.App
	err := c.db.All(&apps)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if err = stream.Send(&app); err != nil {
			return err
		}
	}

	return nil
}

func (c *commandServer) RemoveApp(ctx context.Context, app *pb.App) (*pb.Empty, error) {
	err := c.db.Remove(app)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *commandServer) CreateRelease(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	err := c.db.Save(release)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (c *commandServer) GetRelease(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	err := c.db.Find("id", release.Id, release)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (c *commandServer) UpdateRelease(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	err := c.db.Save(release)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (c *commandServer) UploadRelease(upload pb.Command_UploadReleaseServer) error {
	var releaseID uint64
	var total int64
	var receivedBytes int64
	var err error
	var hash string
	var chunck *pb.Chunck
	var moved bool

	filename := uuid.NewV4().String()
	path := fmt.Sprintf("/tmp/%s", filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
		// if there is an error, the tmp file should be cleaned up
		if err != nil {
			if moved {
				path = fmt.Sprintf("/bundles/%s", hash)
			}

			err = os.Remove(path)
			if err != nil {
				log.Println(err.Error())
				log.Printf("hash '%s' value related to above \n", hash)
			}
		}
	}()

	for {
		chunck, err = upload.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		header := chunck.GetHeader()
		body := chunck.GetBody()

		if header != nil {
			if releaseID != 0 || total != 0 {
				err = fmt.Errorf("chunck header sent multiple times")
				return err
			}
			releaseID = header.ReleaseId
			total = header.Total
		} else if body != nil {
			receivedBytes += int64(len(body.Data))
			file.Write(body.Data)
		}
	}

	if total == 0 || releaseID == 0 {
		err = fmt.Errorf("header is not sent")
		return err
	}

	if receivedBytes != total {
		err = fmt.Errorf("the total amount is not matched")
		return err
	}

	// calculate the hash value
	hash, err = helper.HashFile(path)
	if err != nil {
		return err
	}

	release := pb.Release{}

	err = c.db.One("id", releaseID, &release)
	if err != nil {
		return err
	}

	// initialize buckets
	err = c.db.Init(&pb.App{})
	if err != nil {
		return err
	}

	err = c.db.Init(&pb.Release{})
	if err != nil {
		return err
	}

	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	release.Bundle = hash
	err = tx.Save(release)
	if err != nil {
		return err
	}

	// move the file to bundles folder
	err = os.Rename(path, fmt.Sprintf("/bundles/%s", hash))
	if err != nil {
		return err
	}

	// this is only to clean up the file either from tmp or bundles folder
	moved = true

	err = tx.Commit()
	return err
}

type queryServer struct {
	db *storm.DB
}

func (q *queryServer) GetUpgrade(context.Context, *pb.Upgrade) (*pb.Release, error) {
	return nil, nil
}

func (q *queryServer) DownloadRelease(release *pb.Release, query pb.Query_DownloadReleaseServer) error {
	return nil
}

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

	grpcCommandConfig, err := newGrpcConfig(commandEnv.CA, commandEnv.Crt, commandEnv.Key)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcQueryConfig, err := newGrpcConfig(queryEnv.CA, queryEnv.Crt, queryEnv.Key)
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcCommandServer, err := grpcCommandConfig.createServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcQueryServer, err := grpcQueryConfig.createServer()
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
