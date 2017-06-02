package main

import (
	"fmt"
	"io"
	"os"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
	uuid "github.com/satori/go.uuid"
)

type commandServer struct {
	db *storm.DB
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
	var chunck *pb.Chunck
	var err error
	var moved bool
	var hash string
	var newRelease *pb.Release
	var releases []*pb.Release

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

	// need to convert Upgrades to releases
	if len(header.Upgrades) > 0 {
		releases, err = c.getReleases(newRelease, header.Upgrades)
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
		if moved {

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
	err = os.Rename(path, bundlePath(newRelease))
	if err != nil {
		return err
	}
	// mark the file has been moved to the new location
	// so if the error happens, we can delete the file in
	// the right path
	moved = true

	// connects new release to previous releases
	for _, release := range releases {
		release.NextReleaseId = newRelease.Id
		err = tx.Save(release)
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
