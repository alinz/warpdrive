package data

import (
	"fmt"
	"time"

	"github.com/blang/semver"

	"upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
)

type Release struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	CycleID   int64     `db:"cycle_id,omitempty" json:"-"`
	Platform  Platform  `db:"platform,omitempty" json:"platform"`
	Version   string    `db:"version,omitempty" json:"version"`
	Major     int64     `db:"major,omitempty"`
	Minor     int64     `db:"minor,omitempty"`
	Patch     int64     `db:"patch,omitempty"`
	Build     string    `db:"build,omitempty"`
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

func (r *Release) Load(session db.Database) error {
	if session == nil {
		session = dbSession
	}
	return r.Query(session, db.Cond{"id": r.ID}).One(r)
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

func FindReleases(cycleID int64, platform Platform, version string, note string) ([]*Release, error) {
	var sql string

	if version != "" {
		sql = fmt.Sprintf(`
		SELECT * FROM releases 
		WHERE cycle_id=%d AND version='%s' AND platform=%d AND note like '%%%s%%'
		ORDER BY major, minor, patch DESC, build DESC NULLS FIRST	
	`, cycleID, version, PlatformToInt(platform), note)
	} else {
		sql = fmt.Sprintf(`
		SELECT * FROM releases 
		WHERE cycle_id=%d AND platform=%d AND note like '%%%s%%'
		ORDER BY major, minor, patch DESC, build DESC NULLS FIRST	
	`, cycleID, PlatformToInt(platform), note)
	}

	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil, err
	}

	var releases []*Release
	iter := sqlbuilder.NewIterator(rows)
	iter.All(&releases)

	return releases, nil
}

func FindReleaseByID(cycleID, releaseID int64) (*Release, error) {
	var release Release

	err := release.Find(dbSession, db.Cond{"id": releaseID, "cycle_id": cycleID})
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func FindLockedReleaseByID(cycleID, releaseID int64) (*Release, error) {
	var release Release

	err := release.Find(dbSession, db.Cond{"id": releaseID, "cycle_id": cycleID, "locked": true})
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func FindLockedReleaseByVersion(cycleID int64, platform Platform, version semver.Version) (*Release, error) {
	sql := fmt.Sprintf(`
		SELECT * FROM releases 
		WHERE cycle_id=%d AND platform=%d AND locked=TRUE AND version='%s'
		ORDER BY major, minor, patch DESC, build DESC NULLS FIRST					
	`, cycleID, platform.ValueAsInt(), version.String())

	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil, err
	}

	release := Release{}
	iter := sqlbuilder.NewIterator(rows)
	err = iter.One(&release)

	if err != nil {
		return nil, err
	}

	return &release, nil
}

func FindLatestSoftRelease(cycleID int64, platform Platform, version semver.Version) (*Release, error) {
	// sql := fmt.Sprintf(`
	// 	SELECT * FROM releases
	// 	WHERE cycle_id=%d AND platform=%d AND locked=TRUE AND major=%d AND (minor > %d OR patch > %d)
	// 	ORDER BY major, minor, patch DESC, build DESC NULLS FIRST
	// `, cycleID, platform.ValueAsInt(), version.Major, version.Minor, version.Patch)

	// we need to remove Minor and Patch from the sql statement
	// because of rollback

	sql := fmt.Sprintf(`
		SELECT * FROM releases 
		WHERE cycle_id=%d AND platform=%d AND locked=TRUE AND major=%d
		ORDER BY major, minor, patch DESC, build DESC NULLS FIRST					
	`, cycleID, platform.ValueAsInt(), version.Major)

	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil, err
	}

	release := Release{}
	iter := sqlbuilder.NewIterator(rows)
	err = iter.One(&release)

	if err != nil {
		return nil, err
	}

	foundVersion, err := semver.Make(release.Version)
	if err != nil {
		return nil, err
	}

	if version.Equals(foundVersion) {
		return nil, fmt.Errorf("found same version")
	}

	return &release, nil
}

func FindLatestHardRelease(cycleID int64, platform Platform, version semver.Version) (*Release, error) {
	sql := fmt.Sprintf(`
		SELECT * FROM releases
		WHERE cycle_id=%d AND platform=%d AND locked=TRUE AND major > %d
		ORDER BY major, minor, patch DESC, build DESC NULLS FIRST				
	`, cycleID, platform.ValueAsInt(), version.Major)

	rows, err := dbSession.Query(sql)

	if err != nil {
		return nil, err
	}

	release := Release{}
	iter := sqlbuilder.NewIterator(rows)
	err = iter.One(&release)

	if err != nil {
		return nil, err
	}

	return &release, nil
}
