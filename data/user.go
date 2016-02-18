package data

import (
	"time"

	"github.com/pressly/warpdrive"
	"upper.io/db.v2"
)

type User struct {
	ID        int64     `db:"id,omitempty,pk" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (u *User) LoadByID(id int64) error {
	session := warpdrive.DB
	return session.C("users").Find(db.Cond{
		"id": id,
	}).One(u)
}

func (u *User) LoadByEmailAndPassword(email, password string) error {
	session := warpdrive.DB
	return session.C("users").Find(db.Cond{
		"email":    email,
		"password": password,
	}).One(u)
}

func (u *User) Append(session db.Database) error {
	if session == nil {
		session = warpdrive.DB
	}

	u.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	u.CreatedAt = u.UpdatedAt

	id, err := session.C("users").Append(u)
	if err == nil {
		u.ID = id.(int64)
	}

	return err
}

func (u *User) Update(session db.Database) error {
	if session == nil {
		session = warpdrive.DB
	}

	u.UpdatedAt = time.Now().UTC().Truncate(time.Second)

	return session.
		C("users").
		Find(db.Cond{"id": u.ID}).
		Update(u)
}
