package data

import (
	"fmt"
	"time"

	db "upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
)

type Cycle struct {
	ID         int64     `db:"id,omitempty,pk" json:"id"`
	AppID      int64     `db:"app_id" json:"-"`
	Name       string    `db:"name,omitempty" json:"name"`
	PublicKey  string    `db:"public_key,omitempty" json:"-"`
	PrivateKey string    `db:"private_key,omitempty" json:"-"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func (c Cycle) CollectionName() string {
	return "cycles"
}

func (c Cycle) Query(session db.Database, query db.Cond) db.Result {
	return session.Collection(c.CollectionName()).Find(query)
}

func (c *Cycle) Find(session db.Database, query db.Cond) error {
	return c.Query(session, query).One(c)
}

func (c *Cycle) Load(session db.Database) error {
	if session == nil {
		session = dbSession
	}
	return c.Query(session, db.Cond{"id": c.ID}).One(c)
}

func (c *Cycle) Save(session db.Database) error {
	collection := session.Collection(c.CollectionName())
	var err error

	if c.ID == 0 {
		var id interface{}
		c.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		c.CreatedAt = c.UpdatedAt

		id, err = collection.Insert(c)
		if err == nil {
			c.ID = id.(int64)
		}
	} else {
		c.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		err = collection.
			Find(db.Cond{"id": c.ID}).
			Update(c)
	}

	return err
}

func (c *Cycle) Remove(session db.Database) error {
	return c.Query(session, db.Cond{"id": c.ID}).Delete()
}

func FindCyclesApp(userID, appID int64, name string) []*Cycle {
	app := FindAppByUserIDAppID(userID, appID)

	if app == nil {
		return nil
	}

	sql := fmt.Sprintf(`
		SELECT *
		FROM cycles
		WHERE cycles.app_id=%d AND name='%%%s%%'`, appID, name)
	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil
	}

	var cycles []*Cycle
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&cycles)

	if err != nil {
		return nil
	}

	return cycles
}
