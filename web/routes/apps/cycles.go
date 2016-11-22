package apps

import (
	"net/http"

	"fmt"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func cyclesAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userId").(int64)
	appID, err := web.ParamAsInt64(r, "appId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	name := query.Get("name")

	cycles := services.SearchAppCycles(userID, appID, name)

	web.Respond(w, http.StatusOK, cycles)
}

func getCycleAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userId").(int64)
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

	cycle := services.FindCycleByID(userID, appID, cycleID)

	if cycle == nil {
		web.Respond(w, http.StatusNotFound, fmt.Errorf("cycle not found"))
		return
	}

	web.Respond(w, http.StatusOK, cycle)
}

func createCycleAppHandler(w http.ResponseWriter, r *http.Request) {
}

func getKeyCycleAppHandler(w http.ResponseWriter, r *http.Request) {
}

func updateCycleAppHandler(w http.ResponseWriter, r *http.Request) {
}

func removeCycleAppHandler(w http.ResponseWriter, r *http.Request) {
}
