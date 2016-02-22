package data

import "upper.io/db.v2"

type Data interface {
	CollectionName() string
	Query(db.Database, db.Cond) db.Result
	Find(db.Database, db.Cond) error
	Save(db.Database) error
	Remove(db.Database) error
}
