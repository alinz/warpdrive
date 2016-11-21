package services

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"fmt"

	"github.com/pressly/warpdrive/data"
	"upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
)

// FindUserByEmail try to find a single user by an email. email has to be matached. no
// partial email is permitted.
func FindUserByEmail(email string) *data.User {
	var user data.User

	if err := user.Find(nil, db.Cond{"email": email}); err != nil {
		log.Println(err.Error())
		return nil
	}

	return &user
}

// QueryUsersByEmail this method returns users based on partial email search
func QueryUsersByEmail(name, email string) []*data.User {
	users := data.QueryUsersByEmail(name, email)

	if users == nil {
		users = make([]*data.User, 0)
	}

	return users
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
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), 0)

	if err != nil {
		return nil, fmt.Errorf("couldn't create a proper password")
	}

	user := &data.User{
		Name:     name,
		Email:    email,
		Password: string(hashpass),
	}

	err = user.Save(nil)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(userID int64, name, email, password *string) (*data.User, error) {
	user := FindUserByID(userID)

	if user == nil {
		return nil, fmt.Errorf("you don't have access to this user")
	}

	if name != nil {
		user.Name = *name
	}

	if email != nil {
		user.Email = *email
	}

	if password != nil {
		hashpass, err := bcrypt.GenerateFromPassword([]byte(*password), 0)
		if err != nil {
			return nil, fmt.Errorf("couldn't create a proper password")
		}
		user.Password = string(hashpass)
	}

	err := data.Transaction(func(session sqlbuilder.Tx) error {
		return user.Save(session)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUsersWithinApp(userID, appID int64, name, email string) []*data.User {
	users := data.FindUsersWithinApp(userID, appID, name, email)

	if users == nil {
		users = make([]*data.User, 0)
	}

	return users
}

func AssignUserToApp(currentUserID, userID, appID int64) error {
	app := data.FindAppByUserIDAppID(currentUserID, appID)

	if app == nil {
		return fmt.Errorf("app is not found")
	}

	return data.Transaction(func(session sqlbuilder.Tx) error {
		permission := &data.Permission{
			UserID: userID,
			AppID:  appID,
		}
		return permission.Save(session)
	})
}

func UnassignUserFromApp(currentUserID, userID, appID int64) error {
	app := data.FindAppByUserIDAppID(currentUserID, appID)

	if app == nil {
		return fmt.Errorf("app is not found")
	}

	return data.Transaction(func(session sqlbuilder.Tx) error {
		permission := &data.Permission{}
		err := permission.Find(session, db.Cond{"user_id": userID, "app_id": appID})

		if err != nil {
			return err
		}

		return permission.Remove(session)
	})
}
