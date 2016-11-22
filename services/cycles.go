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

func FindCycleByID(userID, appID, cycleID int64) *data.Cycle {
	app := data.FindAppByUserIDAppID(userID, appID)
	if app == nil {
		return nil
	}

	cycle := &data.Cycle{
		ID: cycleID,
	}

	err := cycle.Load(nil)
	if err != nil {
		return nil
	}

	return cycle
}

func CreateCycle(userID, appID int64, name string) *data.Cycle {
	app := data.FindAppByUserIDAppID(userID, appID)
	if app == nil {
		return nil
	}

	keySize := warpdrive.Conf.Security.KeySize
	privateKey, publicKey, err := crypto.RSAKeyPair(keySize)

	if err != nil {
		return nil
	}

	cycle := &data.Cycle{
		AppID:      appID,
		Name:       name,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	err = cycle.Save(nil)
	
	if err != nil {
		return nil
	}

	return cycle
}
