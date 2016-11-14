package services

import (
	"log"

	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
)

func FindUserByEmailPassword(email, password string) *data.User {
	var user data.User

	if err := user.Find(nil, db.Cond{"email": email, "password": password}); err != nil {
		log.Println(err.Error())
		return nil
	}

	// err := data.Transaction(func(session sqlbuilder.Tx) error {
	// 	user.Query(session, db.Cond{"email": email, "password": password})
	// 	return nil
	// })

	return nil
}
