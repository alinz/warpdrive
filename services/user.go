package services

import (
	"log"

	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
)

// FindUserByEmail try to find a single user by an email. email has to be matached. no
// partial email is permitted.
func FindUserByEmail(email string) *data.User {
	var user data.User

	if err := user.Find(nil, db.Cond{"email": email}); err != nil {
		log.Println(err.Error())
		return nil
	}

	// err := data.Transaction(func(session sqlbuilder.Tx) error {
	// 	user.Query(session, db.Cond{"email": email, "password": password})
	// 	return nil
	// })

	return &user
}

// QueryUsersByEmail this method returns users based on partial email search
func QueryUsersByEmail(email string) []*data.User {
	return data.QueryUsersByEmail(email)
}

// FindUserByID load user by id
func FindUserByID(id int64) *data.User {
	var user data.User
	user.ID = id
	err := user.Load(nil)

	if err != nil {
		return nil
	}

	return &user
}

// CreateUser creates a new user
func CreateUser(name, email, password string) (*data.User, error) {
	user := &data.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	err := user.Save(nil)
	if err != nil {
		return nil, err
	}

	return user, nil
}
