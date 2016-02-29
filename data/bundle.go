package data

import (
	"time"

	"upper.io/db.v2"
)

type Bundle struct {
	ID        int64     `db:"id,omitempty,pk" json:"-"`
	ReleaseID int64     `db:"release_id,omitempty" json:"-"`
	Hash      string    `db:"hash,omitempty" json:"hash"`
	Name      string    `db:"name,omitempty" json:"name"`
	Type      FileType  `db:"type,omitempty" json:"type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (b Bundle) CollectionName() string {
	return "bundles"
}

func (b Bundle) Query(session db.Database, query db.Cond) db.Result {
	return session.C(b.CollectionName()).Find(query)
}

func (b *Bundle) Find(session db.Database, query db.Cond) error {
	return b.Query(session, query).One(b)
}

func (b *Bundle) Save(session db.Database) error {
	collection := session.C(b.CollectionName())
	var err error

	var id interface{}
	b.CreatedAt = time.Now().UTC().Truncate(time.Second)

	id, err = collection.Append(b)
	if err == nil {
		b.ID = id.(int64)
	}

	return err
}

func (b *Bundle) Remove(session db.Database) error {
	return b.Query(session, db.Cond{"id": b.ID}).Remove()
}

func AllBundlesByReleaseID(
	session db.Database,
	releaseID int64,
) ([]*Bundle, error) {
	var bundles []*Bundle

	err := session.
		C("bundles").
		Find(db.Cond{"releaseId": releaseID}).
		All(&bundles)
	if err != nil {
		return nil, err
	}

	return bundles, nil
}
