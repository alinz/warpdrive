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
	createCycle := ctx.Value(constant.CtxKeyParsedBody).(*createCycleRequest)

	if err != nil {
		util.RespondError(w, err)
		return
	}

	cycle, err := service.CreateCycle(*createCycle.Name, appID, userID)
	util.AutoDetectResponse(w, cycle, err)
}

func allAppCyclesHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

	userID := util.LoggedInUserID(ctx)
	appID, err := util.ParamValueAsID(ctx, "appId")

	cycles, err := service.AllAppCycles(appID, userID)
	util.AutoDetectResponse(w, cycles, err)
}

func updateAppCycleHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

	userID := util.LoggedInUserID(ctx)
	appID, err := util.ParamValueAsID(ctx, "appId")
	cycleID, err := util.ParamValueAsID(ctx, "cycleId")
	updateCycle := ctx.Value(constant.CtxKeyParsedBody).(*updateCycleRequest)

	err = service.UpdateAppCycle(*updateCycle.Name, appID, cycleID, userID)
	util.AutoDetectResponse(w, nil, err)
}

func downloadAppCycleConfigHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

	userID := util.LoggedInUserID(ctx)
	appID, err := util.ParamValueAsID(ctx, "appId")
	cycleID, err := util.ParamValueAsID(ctx, "cycleId")

	cycle, err := service.FindAppCycle(appID, cycleID, userID)

	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.Respond(w, 200, struct {
		AppID     int64  `json:"app_id"`
		CycleID   int64  `json:"cycle_id"`
		PublicKey string `json:"public_key"`
	}{
		AppID:     appID,
		CycleID:   cycleID,
		PublicKey: cycle.PublicKey,
	})
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
