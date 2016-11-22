package data

import (
	"fmt"
	"time"

	db "upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
)

type Bundle struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	ReleaseID int64     `db:"release_id,omitempty" json:"-"`
	Hash      string    `db:"hash,omitempty" json:"hash"`
	Name      string    `db:"name,omitempty" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (b Bundle) CollectionName() string {
	return "bundles"
}

func (b Bundle) Query(session db.Database, query db.Cond) db.Result {
	return session.Collection(b.CollectionName()).Find(query)
}

func (b *Bundle) Find(session db.Database, query db.Cond) error {
	return b.Query(session, query).One(b)
}

func (b *Bundle) Load(session db.Database) error {
	if session == nil {
		session = dbSession
	}
	return b.Query(session, db.Cond{"id": b.ID}).One(b)
}

func (b *Bundle) Save(session db.Database) error {
	if session == nil {
		session = dbSession
	}

	collection := session.Collection(b.CollectionName())
	var err error

	if b.ID == 0 {
		var id interface{}
		b.CreatedAt = time.Now().UTC().Truncate(time.Second)

		id, err = collection.Insert(b)
		if err == nil {
			b.ID = id.(int64)
		}
	} else {
		err = collection.
			Find(db.Cond{"id": b.ID}).
			Update(b)
	}

	return err
}

func (b *Bundle) Remove(session db.Database) error {
	return b.Query(session, db.Cond{"id": b.ID}).Delete()
}

func SearchBundlesByName(releaseId int64, name string) ([]*Bundle, error) {
	sql := fmt.Sprintf(`
    SELECT * FROM bundles
    WHERE release_id=%d AND name like '%%%s%%'
  `, releaseId, name)
	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil, nil
	}

	var bundles []*Bundle
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&bundles)

	if err != nil {
		return nil, err
	}

	return bundles, nil
}
