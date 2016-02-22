package apps

import (
	"net/http"

	"github.com/pressly/warpdrive/service"
	"github.com/pressly/warpdrive/web/constant"
	"github.com/pressly/warpdrive/web/util"

	"golang.org/x/net/context"
)

func createAppCycleHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

	userID := util.LoggedInUserID(ctx)
	appID, err := util.ParamValueAsID(ctx, "appId")
	cycleRequest := ctx.Value(constant.CtxKeyParsedBody).(*createCycleRequest)

	if err != nil {
		util.RespondError(w, err)
		return
	}

	cycle, err := service.CreateCycle(*cycleRequest.Name, appID, userID)
	util.AutoDetectResponse(w, cycle, err)
}

func allAppCyclesHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func updateAppCycleHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func downloadAppCycleConfigHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func createAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func updateAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func uploadAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func lockAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func checkVersionAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}

func downloadAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

}
