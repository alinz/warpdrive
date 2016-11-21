package services

import "github.com/pressly/warpdrive/data"

// SearchApps this method returns apps based on partial name search
func SearchApps(userID int64, name string) []*data.App {
	apps := data.SearchAppsByName(userID, name)

	if apps == nil {
		apps = make([]*data.App, 0)
	}

	return apps
}

func FindAppByID(userID, appID int64) *data.App {
	return data.FindAppByUserIDAppID(userID, appID)
}
