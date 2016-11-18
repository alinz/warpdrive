package services

import (
	"log"

	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
)

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

func QueryUsersByEmail(email string) []*data.User {
	return data.QueryUsersByEmail(email)
}
