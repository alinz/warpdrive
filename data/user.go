package data

import (
	"time"

	"upper.io/db.v2"
)

type User struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	Name      string    `db:"name,omitempty" json:"name"`
	Email     string    `db:"email,omitempty" json:"email"`
	Password  string    `db:"password,omitempty" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (u User) CollectionName() string {
	return "users"
}

func (u User) Query(session db.Database, query db.Cond) db.Result {
	return session.C(u.CollectionName()).Find(query)
}

func (u *User) Find(session db.Database, query db.Cond) error {
	return u.Query(session, query).One(u)
}

func (u *User) Save(session db.Database) error {
	collection := session.C(u.CollectionName())
	var err error

	if u.ID == 0 {
		var id interface{}
		u.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		u.CreatedAt = u.UpdatedAt

		id, err = collection.Append(u)
		if err == nil {
			u.ID = id.(int64)
		}
	} else {
		u.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		err = collection.
			Find(db.Cond{"id": u.ID}).
			Update(u)
	}

	return err
}

func (u *User) Remove(session db.Database) error {
	return u.Query(session, db.Cond{"id": u.ID}).Remove()
}
