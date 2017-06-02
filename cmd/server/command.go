package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
	uuid "github.com/satori/go.uuid"
	context "golang.org/x/net/context"
)

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
