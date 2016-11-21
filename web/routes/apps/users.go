package apps

import (
	"net/http"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func usersAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userId").(int64)
	appID, err := web.ParamAsInt64(r, "appId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	name := query.Get("name")
	email := query.Get("email")

	users := services.FindUsersWithinApp(userID, appID, name, email)

	web.Respond(w, http.StatusOK, users)
}

func assignUserToAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUserID := ctx.Value("userId").(int64)

	appID, err := web.ParamAsInt64(r, "appId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	userID, err := web.ParamAsInt64(r, "userId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	err = services.AssignUserToApp(currentUserID, userID, appID)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusCreated, nil)
}

func unassignUserFromAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUserID := ctx.Value("userId").(int64)

	appID, err := web.ParamAsInt64(r, "appId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	userID, err := web.ParamAsInt64(r, "userId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	err = services.UnassignUserFromApp(currentUserID, userID, appID)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusCreated, nil)
}
