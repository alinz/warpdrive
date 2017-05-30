package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	context "golang.org/x/net/context"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/kelseyhightower/envconfig"
	uuid "github.com/satori/go.uuid"

	"github.com/asdine/storm/q"
	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
)

func openDB(path string) (*storm.DB, error) {
	db, err := storm.Open(path, storm.Codec(protobuf.Codec))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func bundlePath(release *pb.Release) string {
	return fmt.Sprintf("/bundles/%s", release.Bundle)
}

type commandServer struct {
	db *storm.DB
}

func (c *commandServer) CreateRelease(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	err := c.db.Save(release)
	if err != nil {
		return nil, err
	}

	return release, nil
}

// if Release.Id is provided, then only the matched one returns.
// if Release.App is provided, then it returns all the releases for that app
// Release.Rollout is also need to be provided
func (c *commandServer) GetRelease(release *pb.Release, stream pb.Command_GetReleaseServer) error {
	if release.Id != 0 {
		err := c.db.One("id", release.Id, release)
		if err != nil {
			return err
		}

		stream.Send(release)
	} else {
		c.db.Select(q.Eq("App", release.App), q.Eq("RolloutAt", release.RolloutAt)).Each(new(pb.Release), func(record interface{}) error {
			stream.Send(record.(*pb.Release))
			return nil
		})
	}

	return nil
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
	//var total int64
	//var receivedBytes int64
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
		log.Println("closing file")
		log.Println(err.Error())
		file.Close()
		// if there is an error, the tmp file should be cleaned up
		if err != nil {

			// it means that the file has already moved to bundle so clean
			// that one instead
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
			fmt.Println("received header", header.ReleaseId)
			if releaseID != 0 {
				err = fmt.Errorf("chunck header sent multiple times")
				return err
			}
			releaseID = header.ReleaseId
			//total = header.Total
		} else if body != nil {
			fmt.Println("received body", len(body.Data))
			//receivedBytes += int64(len(body.Data))
			file.Write(body.Data)
		}
	}

	if releaseID == 0 {
		err = fmt.Errorf("header is not sent")
		return err
	}

	// if receivedBytes != total {
	// 	err = fmt.Errorf("the total amount is not matched")
	// 	return err
	// }

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

	log.Println("moved", path, bundlePath(&release))

	// move the file to bundles folder
	err = os.Rename(path, bundlePath(&release))
	if err != nil {
		return err
	}

	// this is only to clean up the file either from tmp or bundles folder
	moved = true

	err = tx.Commit()

	log.Println("before returnin error:", err.Error())
	return err
}

type queryServer struct {
	db *storm.DB
}

func (qs *queryServer) GetUpgrade(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	if release.Id == 0 {
		return nil, fmt.Errorf("release.id is missing")
	}

	if release.App == "" {
		return nil, fmt.Errorf("release.App is missing")
	}

	if release.Platform == pb.Platform_UNKNOWN {
		return nil, fmt.Errorf("release.Platform is missing")
	}

	if release.RolloutAt == "" {
		return nil, fmt.Errorf("release.RolloutAt is missing")
	}

	// we need to load release object so we can access to
	//
	err := qs.db.Select(q.And(
		q.Eq("id", release.Id),
		q.Eq("app", release.App),
		q.Eq("platform", release.Platform),
		q.Eq("rolloutat", release.RolloutAt),
		q.Eq("lock", true),
	)).First(release)
	if err != nil {
		return nil, err
	}

	// we need to find the latest one
	for {
		err = qs.db.Select(q.And(
			q.Eq("id", release.NextReleaseId),
			q.Eq("app", release.App),
			q.Eq("platform", release.Platform),
			q.Eq("rolloutat", release.RolloutAt),
			q.Eq("lock", true),
		)).First(release)
		if err != nil {
			return nil, err
		}

		if release.NextReleaseId == 0 {
			break
		}
	}

	return release, nil
}

func (qs *queryServer) DownloadRelease(release *pb.Release, stream pb.Query_DownloadReleaseServer) error {
	err := qs.db.One("id", release.Id, release)
	if err != nil {
		return err
	}

	path := bundlePath(release)
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	// we need to send the header first
	info, err := file.Stat()
	if err != nil {
		return err
	}

	chunck := &pb.Chunck{
		Value: &pb.Chunck_Header_{
			Header: &pb.Chunck_Header{
				ReleaseId: release.Id,
				Total:     info.Size(),
			},
		},
	}

	err = stream.Send(chunck)
	if err != nil {
		return err
	}

	// the buffer is 10kb which means
	// for sending 10mb we need to send 1000 messages
	buffer := make([]byte, 10000)
	var n int

	for {
		n, err = file.Read(buffer)
		if err != io.EOF {
			break
		}

		if n > 0 {
			chunck = &pb.Chunck{
				Value: &pb.Chunck_Body_{
					Body: &pb.Chunck_Body{
						Data: buffer,
					},
				},
			}

			err = stream.Send(chunck)
			if err != nil {
				return err
			}
		}
	}

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
