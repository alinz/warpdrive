package service

import (
	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
)

func CreateApp(name string, userID int64) error {
	//crate a tansaction because we need to create an app and
	//assign it to user
	fn := func(session db.Database) error {
		var err error
		//create an app
		app := data.App{Name: name}
		err = app.Save(session)

		if err != nil {
			return err
		}

		//assign agent permision to this app and user
		permission := data.Permission{
			UserID:     userID,
			AppID:      app.ID,
			Permission: data.AGENT,
		}

		err = permission.Save(session)

		return err
	}

	return data.Transaction(fn)
}
