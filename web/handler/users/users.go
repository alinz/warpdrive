package users

import (
	"net/http"

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
