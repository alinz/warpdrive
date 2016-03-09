package data

import (
	"strings"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/config"
	"upper.io/db.v2"
	"upper.io/db.v2/postgresql"
)

func InitDbWithConfig(conf *config.Config) (db.Database, error) {
	var settings = postgresql.ConnectionURL{
		Database: conf.DB.Database,
		Host:     strings.Join(conf.DB.Hosts, ","),
		User:     conf.DB.Username,
		Password: conf.DB.Password,
	}

	db, err := db.Open(postgresql.Adapter, settings)

	return db, err
}

func Transaction(scope func(tx db.Database) error) error {
	tx, err := warpdrive.DB.Transaction()
	if err != nil {
		return err
	}

	defer tx.Close()

	err = scope(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
