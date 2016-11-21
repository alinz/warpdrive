package data

import db "upper.io/db.v2"

//Permission this is reppresentation of Permissions tbale
type Permission struct {
	ID     int64 `db:"id,omitempty,pk" json:"-"`
	UserID int64 `db:"user_id" json:"-"`
	AppID  int64 `db:"app_id" json:"-"`
}

func (p Permission) CollectionName() string {
	return "permissions"
}

func (p Permission) Query(session db.Database, query db.Cond) db.Result {
	return session.Collection(p.CollectionName()).Find(query)
}

func (p *Permission) Find(session db.Database, query db.Cond) error {
	return p.Query(session, query).One(p)
}

func (p *Permission) Save(session db.Database) error {
	collection := session.Collection(p.CollectionName())
	var err error

	if p.ID == 0 {
		var id interface{}
		id, err = collection.Insert(p)
		if err == nil {
			p.ID = id.(int64)
		}
	} else {
		err = collection.
			Find(db.Cond{"id": p.ID}).
			Update(p)
	}

	return err
}

func (p *Permission) Remove(session db.Database) error {
	return p.Query(session, db.Cond{"id": p.ID}).Delete()
}
