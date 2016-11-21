package services

import "github.com/pressly/warpdrive/data"

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
