package service

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"upper.io/db.v2"

	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/web/security"
)

func verifyUserPassword(hashedPassword, password, pepper string) bool {
	if hashedPassword == "" {
		return false
	}
	if password == "" {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+pepper)); err != nil {
		return false
	}
	return true
}

func FindUserByEmailPassword(email, password string) (*data.User, error) {
	user := data.User{}

	if err := user.Find(warpdrive.DB, db.Cond{
		"email": email,
	}); err != nil {
		return nil, errors.New("username or password is incorrect")
	}

	if !verifyUserPassword(user.Password, password, "") {
		return nil, errors.New("username or password is incorrect")
	}

	return &user, nil
}

func FindUserByID(id int64) (*data.User, error) {
	user := data.User{}
	err := user.Find(warpdrive.DB, db.Cond{"id": id})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserByJWT(token *jwt.Token) (*data.User, error) {
	userIDStr := token.Claims["user_id"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)

	if err != nil {
		return nil, err
	}

	return FindUserByID(userID)
}

func GenerateJWT(user *data.User) (string, error) {
	claims := map[string]interface{}{"user_id": fmt.Sprintf("%v", user.ID)}
	tokenStr, err := security.JwtEncode(claims)
	return tokenStr, err
}

func CreateUser(name, email, password string) (*data.User, error) {
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), 0)

	if err != nil {
		return nil, err
	}

	user := data.User{
		Name:     name,
		Email:    email,
		Password: string(hashpass),
	}

	err = user.Save(warpdrive.DB)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func DeleteUserByID(id int64) error {
	user := data.User{ID: id}
	return user.Remove(warpdrive.DB)
}

func UpdateUser(user *data.User) error {
	return user.Save(warpdrive.DB)
}
