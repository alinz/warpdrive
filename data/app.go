package data

import (
	"fmt"
	"time"

	db "upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
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

// SearchAppsByName returns the list of find apps under that user's permission
func SearchAppsByName(userID int64, name string) []*App {
	sql := fmt.Sprintf(`
		SELECT apps.id, apps.name, apps.updated_at, apps.created_at
		FROM apps
		JOIN permissions
		ON apps.id=permissions.app_id
		WHERE permissions.user_id=%d AND apps.name LIKE '%%%s%%'`, userID, name)
	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil
	}

	var apps []*App
	iter := sqlbuilder.NewIterator(rows)
	iter.All(&apps)

	return apps
}
