package data

import (
	"github.com/pressly/warpdrive"

	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/postgresql"
)

//dbSession is private and only access by data package.
//this make sure that services only touch database via data or Transaction.
var dbSession sqlbuilder.Database

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
	dbSession = session

	return session, err
}

//Transaction creates a transaction. I don't want any one elese outside of data package
//access directly to dbSession.
func Transaction(fn func(sqlbuilder.Tx) error) error {
	return dbSession.Tx(fn)
}
