package data

import (
	"github.com/pressly/warpdrive"

	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/postgresql"
)

//DB is global database variable
var DB sqlbuilder.Database

//NewDatabase creates a new database based on what set in global Conf.
//it is better to call this func once and inside your main func.
func NewDatabase() (sqlbuilder.Database, error) {
	conf := warpdrive.Conf
	var settings = postgresql.ConnectionURL{
		Database: conf.DB.Database,
		Host:     conf.DB.Hosts,
		User:     conf.DB.Username,
		Password: conf.DB.Password,
	}

	session, err := postgresql.Open(settings)
	DB = session

	return session, err
}
