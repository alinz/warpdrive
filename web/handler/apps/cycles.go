package apps

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/pressly/warpdrive"
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

	userID := util.LoggedInUserID(ctx)
	appID, _ := util.ParamValueAsID(ctx, "appId")
	cycleID, _ := util.ParamValueAsID(ctx, "cycleId")

	qs := r.URL.Query()

	version := qs.Get("version")
	platform := qs.Get("platform")
	note := qs.Get("note")

	if err := r.ParseMultipartForm(warpdrive.Config.FileUpload.FileMaxSize); err != nil {
		util.RespondError(w, err)
		return
	}

	var filepaths []string
	var filenames []string

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			//we create a clouser here so we can use defer to close the file
			//in case of any errors.
			err := func(fileHeader *multipart.FileHeader) error {
				file, _ := fileHeader.Open()
				defer file.Close()
				path := warpdrive.Config.FileUpload.TempFolder + util.UUID()
				if err := util.CopyDataToFile(file, path); err != nil {
					return err
				}

				filepaths = append(filepaths, path)
				filenames = append(filenames, fileHeader.Filename)

				return nil
			}(fileHeader)

			if err != nil {
				util.RespondError(w, err)
				return
			}
		}
	}

	release, err := service.CreateRelease(
		appID,
		cycleID,
		userID,
		platform,
		version,
		note,
		filenames,
		filepaths,
	)

	util.AutoDetectResponse(w, release, err)
}

func updateAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {
	//will be implemented. not important right now
}

func allAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

	userID := util.LoggedInUserID(ctx)
	appID, _ := util.ParamValueAsID(ctx, "appId")
	cycleID, _ := util.ParamValueAsID(ctx, "cycleId")

	releases, err := service.AllCycleReleases(appID, cycleID, userID)
	util.AutoDetectResponse(w, releases, err)
}

func lockAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {

	userID := util.LoggedInUserID(ctx)
	appID, _ := util.ParamValueAsID(ctx, "appId")
	cycleID, _ := util.ParamValueAsID(ctx, "cycleId")
	releaseID, _ := util.ParamValueAsID(ctx, "releaseId")

	err := service.LockRelease(appID, cycleID, userID, releaseID)
	util.AutoDetectResponse(w, nil, err)
}

func checkVersionAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {
	util.AutoDetectResponse(w, nil, errors.New("Not Implemented yet"))
}

func downloadAppCycleReleaseHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {
	util.AutoDetectResponse(w, nil, errors.New("Not Implemented yet"))
}
