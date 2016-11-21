package apps

import (
	"net/http"

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
