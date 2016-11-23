package services

import (
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
	"upper.io/db.v2/lib/sqlbuilder"
)

func SearchAppCycles(userID, appID int64, name string) ([]*data.Cycle, error) {
	_, err := FindAppByID(userID, appID)
	if err != nil {
		return nil, err
	}

	cycles, err := data.FindCyclesApp(appID, name)
	if err != nil {
		return nil, err
	}

	if cycles == nil {
		cycles = make([]*data.Cycle, 0)
	}

	return cycles, nil
}

func FindCycleByAppIdCycleId(appID, cycleID int64) (*data.Cycle, error) {
	cycle := &data.Cycle{
		ID: cycleID,
	}

	err := cycle.Load(nil)
	if err != nil {
		return nil, err
	}

	return cycle, nil
}

func FindCycleByID(userID, appID, cycleID int64) (*data.Cycle, error) {
	_, err := FindAppByID(userID, appID)
	if err != nil {
		return nil, err
	}

	return FindCycleByAppIdCycleId(appID, cycleID)
}

func CreateCycle(userID, appID int64, name string) (*data.Cycle, error) {
	_, err := FindAppByID(userID, appID)
	if err != nil {
		return nil, err
	}

	keySize := warpdrive.Conf.Security.KeySize
	privateKey, publicKey, err := crypto.RSAKeyPair(keySize)

	if err != nil {
		return nil, err
	}

	cycle := &data.Cycle{
		AppID:      appID,
		Name:       name,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	err = cycle.Save(nil)

	if err != nil {
		return nil, err
	}

	return cycle, nil
}

func GetAppCyclePublicKey(userID, appID, cycleID int64) (string, error) {
	_, err := FindAppByID(userID, appID)
	if err != nil {
		return "", err
	}

	cycle := &data.Cycle{
		ID: cycleID,
	}

	err = cycle.Load(nil)
	if err != nil {
		return "", err
	}

	return cycle.PublicKey, nil
}

func UpdateCycle(userID, appID, cycleID int64, name string) (*data.Cycle, error) {
	_, err := FindAppByID(userID, appID)
	if err != nil {
		return nil, err
	}

	cycle := &data.Cycle{
		ID: cycleID,
	}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		err = cycle.Load(session)
		if err != nil {
			return err
		}

		cycle.Name = name
		err = cycle.Save(session)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cycle, nil
}

func RemoveCycle(userID, appID, cycleID int64) error {
	_, err := FindAppByID(userID, appID)
	if err != nil {
		return err
	}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		cycle := &data.Cycle{
			ID: cycleID,
		}

		return cycle.Remove(session)
	})

	return err
}
