package apps

import (
	"net/http"

	"fmt"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

type createApp struct {
	Name *string `json:"name,required"`
}

type updateApp struct {
	Name *string `json:"name,required"`
}

func getAppsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userId").(int64)
	query := r.URL.Query()
	name := query.Get("name")

	apps := services.SearchApps(userID, name)

	web.Respond(w, http.StatusOK, apps)
}

func getAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userId").(int64)
	appID, err := web.ParamAsInt64(r, "appId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	app := services.FindAppByID(userID, appID)

	if app == nil {
		web.Respond(w, http.StatusNotFound, fmt.Errorf("app not found"))
		return
	}

	web.Respond(w, http.StatusOK, app)
}

func createAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userId").(int64)
	body := ctx.Value("parsed:body").(*createApp)

	app := services.CreateApp(userID, *body.Name)

	if app == nil {
		web.Respond(w, http.StatusBadRequest, fmt.Errorf("app could not be created"))
		return
	}

	web.Respond(w, http.StatusOK, app)
}

func updateAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userId").(int64)
	body := ctx.Value("parsed:body").(*updateApp)
	appID, err := web.ParamAsInt64(r, "appId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	app := services.UpdateApp(userID, appID, *body.Name)

	if app == nil {
		web.Respond(w, http.StatusNotFound, fmt.Errorf("app not found to be updated"))
		return
	}

	web.Respond(w, http.StatusOK, app)
}
