package users

import (
	"net/http"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("q")

	users := services.QueryUsersByEmail(search)

	web.Respond(w, http.StatusOK, users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := web.ParamAsInt64(r, "userId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	user := services.FindUserByID(userID)

	if user == nil {
		web.Respond(w, http.StatusNotFound, nil)
		return
	}

	web.Respond(w, http.StatusOK, user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
