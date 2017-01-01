package apps

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func getReleasesHandler(w http.ResponseWriter, r *http.Request) {
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

	releases, err := services.SearchReleases(appID, cycleID, platform, version, note)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, releases)
}

func checkVersionHandler(w http.ResponseWriter, r *http.Request) {
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

	version := web.ParamAsString(r, "version")
	platform := web.ParamAsString(r, "platform")

	releases, err := services.LatestRelease(appID, cycleID, version, platform)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, releases)
}

func downloadWithReleaseIDHandler(w http.ResponseWriter, r *http.Request) {
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

	encryptedKey, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	err = services.DownloadRelease(appID, cycleID, releaseID, encryptedKey, w)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
	}
}

func downloadWithVersionHandler(w http.ResponseWriter, r *http.Request) {
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

	version := web.ParamAsString(r, "version")
	platform := web.ParamAsString(r, "platform")

	encryptedKey, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	err = services.DownloadVersion(appID, cycleID, platform, version, encryptedKey, w)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
	}
}
