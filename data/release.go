package data

import (
	"fmt"
	"time"

	db "upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
)

type Release struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	CycleID   int64     `db:"cycle_id,omitempty" json:"-"`
	Platform  Platform  `db:"platform,omitempty" json:"platform"`
	Version   Version   `db:"version,omitempty" json:"version"`
	Note      string    `db:"note,omitempty" json:"note"`
	Locked    bool      `db:"locked" json:"locked"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (r Release) CollectionName() string {
	return "releases"
}

func (r Release) Query(session db.Database, query db.Cond) db.Result {
	return session.Collection(r.CollectionName()).Find(query)
}

func (r *Release) Find(session db.Database, query db.Cond) error {
	return r.Query(session, query).One(r)
}

func (r *Release) Save(session db.Database) error {
	if session == nil {
		session = dbSession
	}

	collection := session.Collection(r.CollectionName())
	var err error

	if r.ID == 0 {
		var id interface{}
		r.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		r.CreatedAt = r.UpdatedAt

		id, err = collection.Insert(r)
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
	return r.Query(session, db.Cond{"id": r.ID}).Delete()
}

func FindAllReleases(session db.Database, query db.Cond) ([]*Release, error) {
	collection := session.Collection("releases")
	var releases []*Release
	err := collection.Find(query).All(&releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func FindReleases(cycleID int64, platform Platform, version Version, note string) ([]*Release, error) {
	sql := fmt.Sprintf(`
		SELECT * FROM releases 
		WHERE cycle_id=%d AND version=%d AND platform=%d AND note like '%%%s%%'
	`, cycleID, VersionToInt(version), PlatformToInt(platform), note)
	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil, err
	}

	var releases []*Release
	iter := sqlbuilder.NewIterator(rows)
	iter.All(&releases)

	return releases, nil
}
