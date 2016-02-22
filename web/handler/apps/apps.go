package apps

import (
	"net/http"

	"github.com/pressly/warpdrive/service"
	"github.com/pressly/warpdrive/web/constant"
	"github.com/pressly/warpdrive/web/util"

	"golang.org/x/net/context"
)

func createAppHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	createApp := ctx.Value(constant.CtxKeyParsedBody).(*createAppRequest)
	userID := util.LoggedInUserID(ctx)

	app, err := service.CreateApp(*createApp.Name, userID)
	if err != nil {
		util.RespondError(w, err)
	}

	util.Respond(w, 200, app)
}

func listAllAppsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func updateAppHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	//not implemented yet
}
