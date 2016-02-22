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
