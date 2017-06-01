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

func mustDefineStringValue(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s must have a value", name)
	}
	return nil
}

type commandServer struct {
	db *storm.DB
}

func (c *commandServer) CreateRelease(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	// we need to check for coupld of conditions
	// 1 - `App`, `version`, `rolloutAt` and `platform` must have a value
	// 2 - combining `App`, `version`, `rolloutAt` and `platform` must be unique
	err := mustDefineStringValue("app", release.App)
	if err != nil {
		return nil, err
	}

	err = mustDefineStringValue("version", release.Version)
	if err != nil {
		return nil, err
	}

	err = mustDefineStringValue("rolloutAt", release.RolloutAt)
	if err != nil {
		return nil, err
	}

	if release.Platform == pb.Platform_UNKNOWN {
		return nil, fmt.Errorf("platform must be defined")
	}

	var releases []pb.Release
	err = c.db.Select(q.And(
		q.Eq("App", release.App),
		q.Eq("Version", release.Version),
		q.Eq("RolloutAt", release.RolloutAt),
		q.Eq("Platform", release.Platform),
	)).Find(&releases)

	if err != nil && err != storm.ErrNotFound {
		return nil, err
	}

	if len(releases) > 0 {
		return nil, fmt.Errorf("version must be changed for this app")
	}

	err = c.db.Save(release)
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
		err := c.db.One("Id", release.Id, release)
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

func (c *commandServer) getReleaseByID(id uint64) (*pb.Release, error) {
	release := &pb.Release{}
	err := c.db.One("Id", id, release)
	if err != nil {
		return nil, err
	}
	return release, nil
}

// getReleases, gets the base release and list of versions, and returns the list of
// Release objects which matched the given release specification.
func (c *commandServer) getReleases(release *pb.Release, versions []string) ([]*pb.Release, error) {
	if len(versions) == 0 {
		return nil, fmt.Errorf("versions are empty")
	}

	releases := make([]*pb.Release, 0)

	for _, version := range versions {
		r := &pb.Release{}
		err := c.db.Select(q.And(
			q.Eq("App", release.App),
			q.Eq("RolloutAt", release.RolloutAt),
			q.Eq("Platform", release.Platform),
			q.Eq("Version", version),
			q.Eq("NextReleaseId", 0),
			q.Eq("Lock", true),
		)).First(r)

		if err != nil {
			return nil, fmt.Errorf("version %s is not compatiable with this release", version)
		}

		releases = append(releases, r)
	}

	return releases, nil
}

func (c *commandServer) UploadRelease(upload pb.Command_UploadReleaseServer) error {
	var releaseID uint64
	//var total int64
	//var receivedBytes int64
	var err error
	var hash string
	var chunck *pb.Chunck
	var moved bool

	var release *pb.Release
	var releases []*pb.Release

	filename := uuid.NewV4().String()
	path := fmt.Sprintf("/bundles/%s", filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
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
			}
		}
	}()

	for {
		chunck, err = upload.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		header := chunck.GetHeader()
		body := chunck.GetBody()

		if header != nil {
			if releaseID != 0 {
				err = fmt.Errorf("chunck header sent multiple times")
				return err
			}
			releaseID = header.ReleaseId
			upgrades := header.Upgrades

			release, err = c.getReleaseByID(releaseID)
			if err != nil {
				return err
			}

			if upgrades != nil && len(upgrades) > 0 {
				// this update is not root, so we need to check all versions exists with the same content
				releases, err = c.getReleases(release, upgrades)
				if err != nil {
					return err
				}
			}
		} else if body != nil {
			file.Write(body.Data)
		}
	}

	if releaseID == 0 {
		err = fmt.Errorf("header is not sent")
		return err
	}

	// calculate the hash value
	hash, err = helper.HashFile(path)
	if err != nil {
		return err
	}

	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	release.Bundle = hash
	release.Lock = true

	err = tx.Save(release)
	if err != nil {
		return err
	}

	// rename the bundle in the same folder.
	// NOTE: if you don't, you will get `invalid cross-device link` error
	err = os.Rename(path, bundlePath(release))
	if err != nil {
		return err
	}

	// this is only to clean up the file either from tmp or bundles folder
	moved = true

	for _, prev := range releases {
		prev.NextReleaseId = release.Id
		err = tx.Save(prev)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err == nil {
		err = upload.SendAndClose(release)
	}

	return err
}

type queryServer struct {
	db *storm.DB
}

func (qs *queryServer) GetUpgrade(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	if release.Version == "" {
		return nil, fmt.Errorf("release.Version is missing")
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
		q.Eq("App", release.App),
		q.Eq("Platform", release.Platform),
		q.Eq("RolloutAt", release.RolloutAt),
		q.Eq("Version", release.Version),
		q.Eq("Lock", true),
	)).First(release)
	if err != nil {
		return nil, err
	}

	// we need to find the latest one
	for {
		err = qs.db.Select(q.And(
			q.Eq("Id", release.NextReleaseId),
			q.Eq("App", release.App),
			q.Eq("Platform", release.Platform),
			q.Eq("RolloutAt", release.RolloutAt),
			q.Eq("Lock", true),
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

	chunck := &pb.Chunck{
		Value: &pb.Chunck_Header_{
			Header: &pb.Chunck_Header{
				ReleaseId: release.Id,
			},
		},
	}

	err = stream.Send(chunck)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1000)
	var n int

	for {
		n, err = file.Read(buffer)
		if err != io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if n > 0 {
			chunck = &pb.Chunck{
				Value: &pb.Chunck_Body_{
					Body: &pb.Chunck_Body{
						Data: buffer[:n],
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
