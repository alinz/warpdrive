package services

import (
	"fmt"

	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
)

var (
	// ErrorLogin returns when email or password are wrong
	ErrorLogin = fmt.Errorf("email and/or password are incorrect")
)

func CreateSession(email, password string) error {
	var user data.User

	if err := user.Find(nil, db.Cond{"email": email, "password": password}); err != nil {
		log.Println(err.Error())
		return ErrorLogin
	}

	// err := data.Transaction(func(session sqlbuilder.Tx) error {
	// 	user.Query(session, db.Cond{"email": email, "password": password})
	// 	return nil
	// })

	return nil
}

func DestorySession() error {
	return nil
}
