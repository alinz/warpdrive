package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/web/security"
)

type userService struct {
	user *data.User
}

func (u *userService) FindUserByEmailPassword(email, password string) error {
	user := data.User{}

	if err := user.LoadByEmailAndPassword(email, password); err != nil {
		return err
	}

	u.user = &user

	return nil
}

func (u *userService) FindUserByID(id int64) error {
	user := data.User{}
	err := user.LoadByID(id)

	if err == nil {
		u.user = &user
	}

	return err
}

func (u *userService) FindUserByJWT(token *jwt.Token) error {
	userIDStr := token.Claims["user_id"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)

	if err != nil {
		return err
	}

	return u.FindUserByID(userID)
}

func (u userService) UserID() int64 {
	var id int64

	if u.user != nil {
		id = u.user.ID
	}

	return id
}

func (u userService) GenerateJWT() (string, error) {
	if u.user == nil {
		return "", errors.New("user is not loaded to generate jwt token")
	}

	claims := map[string]interface{}{"user_id": fmt.Sprintf("%v", u.user.ID)}
	tokenStr, err := security.JwtEncode(claims)
	return tokenStr, err
}
