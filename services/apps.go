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

func FindAppByID(userID, appID int64) (*data.App, error) {
	app := data.FindAppByUserIDAppID(userID, appID)

	if app == nil {
		return nil, ErrAppNotFound
	}

	return app, nil
}

func CreateApp(userID int64, name string) (*data.App, error) {
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
		return nil, err
	}

	return app, nil
}

func UpdateApp(userID, appID int64, name string) (*data.App, error) {
	app, err := FindAppByID(userID, appID)

	if err != nil {
		return nil, err
	}

	app.Name = name

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		return app.Save(session)
	})

	if err != nil {
		return nil, err
	}

	return app, nil
}
