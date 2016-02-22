package service

import (
	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
)

func CreateApp(name string, userID int64) (*data.App, error) {
	app := data.App{Name: name}
	//crate a tansaction because we need to create an app and
	//assign it to user
	fn := func(session db.Database) error {
		var err error
		//create an app
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

		if err != nil {
			return err
		}

		//we also need to assign the permision to root
		if userID == 1 {
			permission = data.Permission{
				UserID:     1,
				AppID:      app.ID,
				Permission: data.AGENT,
			}

			err = permission.Save(session)
		}

		return err
	}

	err := data.Transaction(fn)

	if err != nil {
		return nil, err
	}

	return &app, err
}
