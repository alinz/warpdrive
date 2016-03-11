package service

import (
	"errors"

	"upper.io/db.v2"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
)

func CreateCycle(name string, appID, userID int64) (*data.Cycle, error) {
	if !HasPermissionToCreateAppCycle(appID, userID) {
		return nil, errors.New("No access to this app")
	}

	cycle := data.Cycle{
		Name:  name,
		AppID: appID,
	}

	keySize := warpdrive.Config.Security.KeySize

	fn := func(session db.Database) error {
		privateKey, publicKey, err := crypto.RSAKeyPair(keySize)
		if err != nil {
			return err
		}

		cycle.PrivateKey = privateKey
		cycle.PublicKey = publicKey

		err = cycle.Save(session)
		if err != nil {
			return err
		}

		return err
	}

	err := data.Transaction(fn)

	if err != nil {
		return nil, err
	}

	return &cycle, nil
}

func AllAppCycles(appID, userID int64) ([]*data.Cycle, error) {
	if !HasPermissionToCreateAppCycle(appID, userID) {
		return nil, errors.New("No access to this app")
	}

	builder := warpdrive.DB.Builder()
	q := builder.
		Select("cycles.id",
			"cycles.name",
			"cycles.updated_at",
			"cycles.created_at").
		From("cycles").
		Where("app_id=?", appID)

	var cycles []*data.Cycle
	err := q.Iterator().All(&cycles)

	return cycles, err
}

func UpdateAppCycle(name string, appID, cycleID, userID int64) error {
	if !HasPermissionToCreateAppCycle(appID, userID) {
		return errors.New("No access to this app")
	}

	cycle := data.Cycle{
		ID:    cycleID,
		AppID: appID,
		Name:  name,
	}

	return cycle.Save(warpdrive.DB)
}

func FindAppCycle(appID, cycleID, userID int64) (*data.Cycle, error) {
	if !HasPermissionToCreateAppCycle(appID, userID) {
		return nil, errors.New("No access to this app")
	}

	cycle := data.Cycle{}

	err := cycle.Find(warpdrive.DB, db.Cond{
		"id": cycleID,
	})

	return &cycle, err
}
