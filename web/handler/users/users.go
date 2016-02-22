package users

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/service"
	"github.com/pressly/warpdrive/web/constant"
	"github.com/pressly/warpdrive/web/util"

	"golang.org/x/net/context"
)

func createUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	createUser := ctx.Value(constant.CtxKeyParsedBody).(*createUserRequest)

	_, err := service.CreateUser(*createUser.Name, *createUser.Email, *createUser.Password)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.Respond(w, 200, nil)
}

func deleteUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := util.ParamValueAsID(ctx, "userId")
	if err != nil {
		util.RespondError(w, err)
		return
	}

	err = service.DeleteUserByID(userID)

	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.Respond(w, 200, nil)
}

func updateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	updateUser := ctx.Value(constant.CtxKeyParsedBody).(*updateUserRequest)
	var user *data.User

	userID, err := util.ParamValueAsID(ctx, "userId")
	if err != nil {
		jwt := ctx.Value(constant.CtxJWT).(*jwt.Token)
		user, err = service.FindUserByJWT(jwt)
	} else {
		//userId must be either logged in userId or root user
		if userID != util.LoggedInUserID(ctx) && !util.UserIsRoot(ctx) {
			err = constant.ErrorAuthorizeAccess
		} else {
			user, err = service.FindUserByID(userID)
		}
	}

	if err != nil {
		util.RespondError(w, err)
		return
	}

	user.Name = updateUser.Name
	user.Email = updateUser.Email
	user.Password = updateUser.Password

	err = service.UpdateUser(user)

	if err != nil {
		util.RespondError(w, err)
	} else {
		util.Respond(w, 200, nil)
	}
}
