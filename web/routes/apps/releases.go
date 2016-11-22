package apps

import (
	"net/http"

	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

type createRelease struct {
	Platform *data.Platform `json:"platform,required"`
	Version  *data.Version  `json:"version,required"`
	Note     *string        `json:"note"`
}

func getReleasesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	appID, err := web.ParamAsInt64(r, "appId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	cycleID, err := web.ParamAsInt64(r, "cycleId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	platform := query.Get("platform")
	version := query.Get("version")
	note := query.Get("note")

	releases, err := services.SearchReleases(userID, appID, cycleID, platform, version, note)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, releases)
}

func getReleaseHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	appID, err := web.ParamAsInt64(r, "appId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	cycleID, err := web.ParamAsInt64(r, "cycleId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	releaseID, err := web.ParamAsInt64(r, "releaseId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	release, err := services.FindReleaseByID(userID, appID, cycleID, releaseID)
	if err != nil {
		web.Respond(w, http.StatusNotFound, err)
		return
	}

	web.Respond(w, http.StatusOK, release)
}

func createReleaseHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	body := ctx.Value("parsed:body").(*createRelease)

	appID, err := web.ParamAsInt64(r, "appId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	cycleID, err := web.ParamAsInt64(r, "cycleId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	release, err := services.CreateRelease(userID, appID, cycleID, *body.Platform, *body.Version, *body.Note)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, release)
}

func updateReleaseHandler(w http.ResponseWriter, r *http.Request) {

}

func removeReleaseHandler(w http.ResponseWriter, r *http.Request) {

}

func lockReleaseHandler(w http.ResponseWriter, r *http.Request) {

}

func unlockReleaseHandler(w http.ResponseWriter, r *http.Request) {

}
