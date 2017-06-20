package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
	"github.com/pressly/warpdrive/server/config"
	uuid "github.com/satori/go.uuid"
)

type commandServer struct {
	db   *storm.DB
	conf *config.Config
}

func (c *commandServer) isUnique(release *pb.Release) bool {
	var r pb.Release
	err := c.db.Select(q.And(
		q.Eq("App", release.App),
		q.Eq("RolloutAt", release.RolloutAt),
		q.Eq("Platform", release.Platform),
		q.Eq("Version", release.Version),
	)).First(&r)

	if err == storm.ErrNotFound {
		return true
	}

	if err != nil {
		log.Println("isUnique:", err.Error())
	}

	return false
}

// getRelease find release based on given release template and version number and rootVersion
func (c *commandServer) getRelease(release *pb.Release, version, rootVersion string) (*pb.Release, error) {
	root := &pb.Release{}
	err := c.db.Select(q.And(
		q.Eq("App", release.App),
		q.Eq("RolloutAt", release.RolloutAt),
		q.Eq("Platform", release.Platform),
		q.Eq("Version", rootVersion),
		q.Eq("Lock", true),
	)).First(root)

	if err != nil {
		return nil, fmt.Errorf("root version not found")
	}

	next := root
	for {
		if next.Version == version {
			if next.NextReleaseId == 0 {
				return next, nil
			}
			return nil, fmt.Errorf("previous version already connected to another version")
		}

		if next.NextReleaseId == 0 {
			break
		}

		err = c.db.Select(q.And(
			q.Eq("Id", next.NextReleaseId),
			q.Eq("App", release.App),
			q.Eq("RolloutAt", release.RolloutAt),
			q.Eq("Platform", release.Platform),
			q.Eq("Lock", true),
		)).First(next)

		if err != nil {
			return nil, fmt.Errorf("version %s is not compatiable with this release", version)
		}
	}

	return nil, fmt.Errorf("previous version not found")
}

func (c *commandServer) UploadRelease(upload pb.Command_UploadReleaseServer) error {
	var chunck *pb.Chunck
	var err error
	var moved bool
	var hash string
	var newRelease *pb.Release
	var prevRelease *pb.Release

	chunck, err = upload.Recv()
	if err != nil {
		return err
	}

	header := chunck.GetHeader()
	if header == nil {
		return fmt.Errorf("header is not being sent")
	}

	newRelease = header.Release
	if newRelease == nil {
		return fmt.Errorf("info about new release not found")
	}

	// check if newRelease is unique
	if !c.isUnique(newRelease) {
		return fmt.Errorf("release already exists")
	}

	// need to find previous version
	header.Root = strings.Trim(header.Root, " \t")
	header.Upgrade = strings.Trim(header.Upgrade, " \t")
	if header.Root != "" {
		if header.Upgrade == "" {
			log.Println("root:", header.Root, "upgrade:", header.Upgrade)
			return fmt.Errorf("upgrade is missing since root version is provided")
		}

		prevRelease, err = c.getRelease(newRelease, header.Upgrade, header.Root)
		if err != nil {
			return err
		}
	}

	// file operation
	filename := uuid.NewV4().String()
	path := fmt.Sprintf("/bundles/%s", filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()

		if err != nil {
			if moved {
				os.Remove(filepath.Join(c.conf.BundlePath, newRelease.Bundle))
			} else {
				os.Remove(path)
			}
		}
	}()

	for {
		chunck, err = upload.Recv()
		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return err
		}

		body := chunck.GetBody()
		if body == nil {
			err = fmt.Errorf("something went wrong with uploading data")
			return err
		}

		_, err = file.Write(body.Data)
		if err != nil {
			return err
		}
	}

	// calculate the hash value of uploaded file
	hash, err = helper.HashFile(path)
	if err != nil {
		return err
	}

	// start the transaction
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// the logic of transaction goes here
	/////////////////////////////////////
	newRelease.Bundle = hash
	newRelease.Lock = true

	err = tx.Save(newRelease)
	if err != nil {
		return err
	}

	// moved the temp bundle to the bundle section
	err = os.Rename(path, filepath.Join(c.conf.BundlePath, newRelease.Bundle))
	if err != nil {
		return err
	}

	// mark the file has been moved to the new location
	// so if the error happens, we can delete the file in
	// the right path
	moved = true

	// connects new release to previous release
	if prevRelease != nil {
		prevRelease.NextReleaseId = newRelease.Id
		err = tx.Save(prevRelease)
		if err != nil {
			return err
		}
	}
	/////////////////////////////////////

	err = tx.Commit()
	if err != nil {
		return err
	}

	// this error is not that important that we need to rollback
	return upload.SendAndClose(newRelease)
}

// NewCommandServer creates a Command server
func NewCommandServer(db *storm.DB, conf *config.Config) pb.CommandServer {
	return &commandServer{db, conf}
}
