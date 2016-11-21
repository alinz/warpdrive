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
