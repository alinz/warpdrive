package users

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/warpdrive/service"
	"github.com/pressly/warpdrive/web/constant"
	"github.com/pressly/warpdrive/web/util"

	"golang.org/x/net/context"
)

func createUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	createUser := ctx.Value(constant.CtxKeyParsedBody).(*createUserRequest)

	userService := service.New()

	err := userService.CreateUser(*createUser.Name, *createUser.Email, *createUser.Password)
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

	userService := service.New()
	err = userService.DeleteUserByID(userID)

	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.Respond(w, 200, nil)
}

func updateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	updateUser := ctx.Value(constant.CtxKeyParsedBody).(*updateUserRequest)
	userService := service.New()

	userID, err := util.ParamValueAsID(ctx, "userId")
	if err != nil {
		jwt := ctx.Value(constant.CtxJWT).(*jwt.Token)
		err = userService.FindUserByJWT(jwt)
	} else {
		//userId must be either logged in userId or root user
		if userID != util.LoggedInUserID(ctx) && !util.UserIsRoot(ctx) {
			err = constant.ErrorAuthorizeAccess
		} else {
			err = userService.FindUserByID(userID)
		}
	}

	if err != nil {
		util.RespondError(w, err)
		return
	}

	err = userService.UpdateUser(updateUser.Name, updateUser.Email, updateUser.Password)

	if err != nil {
		util.RespondError(w, err)
	} else {
		util.Respond(w, 200, nil)
	}
}
