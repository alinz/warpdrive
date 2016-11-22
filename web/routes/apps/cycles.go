package apps

import (
	"net/http"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

type createCycle struct {
	Name *string `json:"name,required"`
}

type updateCycle struct {
	Name *string `json:"name,required"`
}

func getCyclesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	appID, err := web.ParamAsInt64(r, "appId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	name := query.Get("name")

	cycles, err := services.SearchAppCycles(userID, appID, name)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, cycles)
}

func getCycleHandler(w http.ResponseWriter, r *http.Request) {
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

	cycle, err := services.FindCycleByID(userID, appID, cycleID)

	if err != nil {
		web.Respond(w, http.StatusNotFound, err)
		return
	}

	web.Respond(w, http.StatusOK, cycle)
}

func createCycleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	appID, err := web.ParamAsInt64(r, "appId")
	body := ctx.Value("parsed:body").(*createCycle)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	cycle, err := services.CreateCycle(userID, appID, *body.Name)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, cycle)
}

func getCycleKeyHandler(w http.ResponseWriter, r *http.Request) {
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

	publicKey, err := services.GetAppCyclePublicKey(userID, appID, cycleID)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, publicKey)
}

func updateCycleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)
	body := ctx.Value("parsed:body").(*updateCycle)
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

	cycle, err := services.UpdateCycle(userID, appID, cycleID, *body.Name)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, cycle)
}

func removeCycleHandler(w http.ResponseWriter, r *http.Request) {
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

	err = services.RemoveCycle(userID, appID, cycleID)

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, nil)
}
