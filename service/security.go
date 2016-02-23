package service

import (
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
)

func HasPermissionToCreateAppCycle(appID, userID int64) bool {
	permission := data.Permission{}
	total, err := permission.Query(warpdrive.DB, db.Cond{
		"app_id":  appID,
		"user_id": userID,
	}).Count()

	return err == nil && total == 1
}

func HashPermissionToAccessCycle(appID, cycleID, userID int64) bool {
	builder := warpdrive.DB.Builder()
	q := builder.
		Select("cycles.id").
		From("cycles").
		Join("apps").
		On("apps.id=cycles.app_id").
		Join("permissions").
		On("apps.id=permissions.app_id").
		Where("permissions.app_id=? AND permissions.user_id=? AND cycles.id=?", appID, userID, cycleID)

	type Result struct {
		ID int64 `db:"id,omitempty,pk"`
	}

	result := Result{}
	err := q.Iterator().One(&result)

	return err == nil && result.ID != 0
}

func HasReleaseLocked(releaseID int64) bool {
	release := data.Release{}
	err := release.Find(warpdrive.DB, db.Cond{"id": releaseID})
	//if there is an error, it means that it has been locked, other wise return
	//whatever the value of lock.
	return err != nil || release.Locked == true
}
