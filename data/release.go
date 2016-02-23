package data

import (
	"time"

	"upper.io/db.v2"
)

type Release struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	CycleID   int64     `db:"cycle_id,omitempty" json:"-"`
	Platform  Platform  `db:"platform,omitempty" json:"platform"`
	Version   Version   `db:"version,omitempty" json:"version"`
	Note      string    `db:"note,omitempty" json:"note"`
	Lock      bool      `db:"lock,omitempty" json:"lock"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (r Release) CollectionName() string {
	return "releases"
}

func (r Release) Query(session db.Database, query db.Cond) db.Result {
	return session.C(r.CollectionName()).Find(query)
}

func (r *Release) Find(session db.Database, query db.Cond) error {
	return r.Query(session, query).One(r)
}

func (r *Release) Save(session db.Database) error {
	collection := session.C(r.CollectionName())
	var err error

	if r.ID == 0 {
		var id interface{}
		r.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		r.CreatedAt = r.UpdatedAt

		id, err = collection.Append(r)
		if err == nil {
			r.ID = id.(int64)
		}
	} else {
		r.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		err = collection.
			Find(db.Cond{"id": r.ID}).
			Update(r)
	}

	return err
}

func (r *Release) Remove(session db.Database) error {
	return r.Query(session, db.Cond{"id": r.ID}).Remove()
}

func FindAllReleases(session db.Database, query db.Cond) ([]*Release, error) {
	collection := session.C("releases")
	var releases []*Release
	err := collection.Find(query).All(&releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}
