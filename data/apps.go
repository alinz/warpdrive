package data

import "upper.io/db.v2"

type Apps []*App

func (a Apps) CollectionName() string {
	return "apps"
}

func (a Apps) Query(session db.Database, query db.Cond) db.Result {
	return session.C(a.CollectionName()).Find(query)
}

func (a *Apps) Find(session db.Database, query db.Cond) error {
	return a.Query(session, query).All(&a)
}

func (a Apps) Save(session db.Database) error {
	var err error
	for _, app := range a {
		err = app.Save(session)
		if err != nil {
			return err
		}
	}
	return err
}

func (a Apps) Remove(session db.Database) error {
	var err error
	for _, app := range a {
		err = app.Remove(session)
		if err != nil {
			return err
		}
	}
	return err
}
