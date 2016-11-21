package services

import (
	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2/lib/sqlbuilder"
)

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

func CreateApp(userID int64, name string) *data.App {
	app := &data.App{
		Name: name,
	}

	err := data.Transaction(func(session sqlbuilder.Tx) error {
		err := app.Save(session)
		if err != nil {
			return err
		}

		permission := &data.Permission{
			AppID:  app.ID,
			UserID: userID,
		}

		err = permission.Save(session)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil
	}

	return app
}

func UpdateApp(userID, appID int64, name string) *data.App {
	app := FindAppByID(userID, appID)

	if app != nil {
		app.Name = name

		err := data.Transaction(func(session sqlbuilder.Tx) error {
			return app.Save(session)
		})

		if err != nil {
			return nil
		}

		return app
	}

	return nil
}
