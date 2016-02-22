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
