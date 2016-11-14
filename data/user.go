package data

import (
	"time"

	db "upper.io/db.v2"
)

type User struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	Name      string    `db:"name,omitempty" json:"name"`
	Email     string    `db:"email,omitempty" json:"email"`
	Password  string    `db:"password,omitempty" json:"-"`
	UpdatedAt time.Time `db:"updated_at,omitempty" json:"updated_at"`
	CreatedAt time.Time `db:"created_at,omitempty" json:"created_at"`
}

func (u User) CollectionName() string {
	return "users"
}

func (u User) Query(session db.Database, query db.Cond) db.Result {
	if session == nil {
		session = dbSession
	}
	return session.Collection(u.CollectionName()).Find(query)
}

func (u *User) Load(session db.Database) error {
	if session == nil {
		session = dbSession
	}
	return u.Query(session, db.Cond{"id": u.ID}).One(u)
}

func (u *User) Find(session db.Database, query db.Cond) error {
	if session == nil {
		session = dbSession
	}
	return u.Query(session, query).One(u)
}

func (u *User) Save(session db.Database) error {
	if session == nil {
		session = dbSession
	}
	collection := session.Collection(u.CollectionName())
	var err error

	u.UpdatedAt = time.Now().UTC().Truncate(time.Second)

	if u.ID == 0 {
		var id interface{}
		u.CreatedAt = u.UpdatedAt

		id, err = collection.Insert(u)
		if err == nil {
			u.ID = id.(int64)
		}
	} else {
		err = collection.
			Find(db.Cond{"id": u.ID}).
			Update(u)
	}

	return err
}

func (u *User) Remove(session db.Database) error {
	if session == nil {
		session = dbSession
	}
	return u.Query(session, db.Cond{"id": u.ID}).Delete()
}
