package apps

import (
	"net/http"

	"fmt"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func getAppsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userid").(int64)
	query := r.URL.Query()
	search := query.Get("q")

	apps := services.SearchApps(userID, search)

	web.Respond(w, http.StatusOK, apps)
}

func getAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userid").(int64)
	appID, err := web.ParamAsInt64(r, "userId")

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
