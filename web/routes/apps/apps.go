package apps

import (
	"net/http"

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
	userID := ctx.Value("user:id").(int64)
	query := r.URL.Query()
	name := query.Get("name")

	apps := services.SearchApps(userID, name)

	web.Respond(w, http.StatusOK, apps)
}

func getAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	appID, err := web.ParamAsInt64(r, "appId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	app, err := services.FindAppByID(userID, appID)

	if err != nil {
		web.Respond(w, http.StatusNotFound, err)
		return
	}

	web.Respond(w, http.StatusOK, app)
}

func createAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	body := ctx.Value("parsed:body").(*createApp)

	app, err := services.CreateApp(userID, *body.Name)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, app)
}

func updateAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	body := ctx.Value("parsed:body").(*updateApp)
	appID, err := web.ParamAsInt64(r, "appId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	app, err := services.UpdateApp(userID, appID, *body.Name)

	if err != nil {
		web.Respond(w, http.StatusNotFound, err)
		return
	}

	web.Respond(w, http.StatusOK, app)
}
