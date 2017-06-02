package main

import (
	"fmt"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
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
