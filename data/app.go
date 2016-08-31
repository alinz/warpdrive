package data

import (
	"time"

	db "upper.io/db.v2"
)

type App struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	Name      string    `db:"name,omitempty" json:"name"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (a App) CollectionName() string {
	return "apps"
}

func (a App) Query(session db.Database, query db.Cond) db.Result {
	return session.Collection(a.CollectionName()).Find(query)
}

func (a *App) Find(session db.Database, query db.Cond) error {
	return a.Query(session, query).One(a)
}

func (a *App) Save(session db.Database) error {
	collection := session.Collection(a.CollectionName())
	var err error

	if a.ID == 0 {
		var id interface{}
		a.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		a.CreatedAt = a.UpdatedAt

		id, err = collection.Insert(a)
		if err == nil {
			a.ID = id.(int64)
		}
	} else {
		a.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		err = collection.
			Find(db.Cond{"id": a.ID}).
			Update(a)
	}

	return err
}

func (a *App) Remove(session db.Database) error {
	return a.Query(session, db.Cond{"id": a.ID}).Delete()
}

type AppWithPermission struct {
	App        `db:",inline"`
	Permission PermissionType `db:"permission" json:"permission"`
}

type Apps []*App

func (a Apps) CollectionName() string {
	return "apps"
}

func (a Apps) Query(session db.Database, query db.Cond) db.Result {
	return session.Collection(a.CollectionName()).Find(query)
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
