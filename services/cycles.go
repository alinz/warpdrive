package services

import (
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
)

func SearchAppCycles(userID, appID int64, name string) []*data.Cycle {
	cycles := data.FindCyclesApp(userID, appID, name)

	if cycles == nil {
		cycles = make([]*data.Cycle, 0)
	}

	return cycles
}

func FindCycleByID(userID, appID, cycleID int64) (*data.Cycle, error) {
	_, err := FindAppByID(userID, appID)
	if err != nil {
		return nil, err
	}

	cycle := &data.Cycle{
		ID: cycleID,
	}

	err = cycle.Load(nil)
	if err != nil {
		return nil, err
	}

	return cycle, nil
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
