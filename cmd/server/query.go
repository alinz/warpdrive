package main

import (
	"fmt"
	"io"
	"os"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	pb "github.com/pressly/warpdrive/proto"
	context "golang.org/x/net/context"
)

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
