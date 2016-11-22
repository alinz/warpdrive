package data

import (
	"fmt"
	"time"

	"strings"

	db "upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
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

func QueryUsersByEmail(name, email string) []*User {
	name = strings.ToLower(name)
	email = strings.ToLower(email)
	sql := fmt.Sprintf(`
		SELECT * from users WHERE lower(email) LIKE '%%%s%%' AND lower(name) LIKE '%%%s%%'
	`, email, name)
	rows, err := dbSession.Query(sql)
	if err != nil {
		return nil
	}

	var users []*User
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&users)

	if err != nil {
		return nil
	}

	return users
}

func FindUsersWithinApp(userID, appID int64, name, email string) []*User {
	// first we need to make sure that userID has access to appID
	app := FindAppByUserIDAppID(userID, appID)

	if app == nil {
		return nil
	}

	name = strings.ToLower(name)
	email = strings.ToLower(email)

	sql := fmt.Sprintf(`
		SELECT users.id, users.name, users.email, users.updated_at, users.created_at
		FROM users
		JOIN permissions
		ON permissions.user_id=users.id
		WHERE permissions.app_id=%d AND lower(users.name) LIKE '%%%s%%' AND lower(users.email) LIKE '%%%s%%'
	`, appID, name, email)
	rows, err := dbSession.Query(sql)
	if err != nil {
		return nil
	}

	var users []*User
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&users)

	if err != nil {
		return nil
	}

	return users
}
