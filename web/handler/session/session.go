package session

import (
	"net/http"

	"github.com/pressly/warpdrive/web/security"

	"github.com/pressly/warpdrive/service"
	"github.com/pressly/warpdrive/web/constant"
	"github.com/pressly/warpdrive/web/util"

	"golang.org/x/net/context"
)

func startSessionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	login := ctx.Value(constant.CtxKeyParsedBody).(*loginRequest)

	user, err := service.FindUserByEmailPassword(*login.Email, *login.Password)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	jwt, err := service.GenerateJWT(user)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	security.SetJwtCookie(jwt, w)

	util.Respond(w, 200, struct {
		JWT string `json:"jwt"`
	}{jwt})
}

func endSessionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	security.RemoveJwtCookie(w)
	util.Respond(w, 200, nil)
}
