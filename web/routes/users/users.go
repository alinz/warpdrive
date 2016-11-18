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
	w.Write([]byte(`¯\_(ツ)_/¯`))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
