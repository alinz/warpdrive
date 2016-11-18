package users

import (
	"net/http"

	"github.com/pressly/chi"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", getUsersHandler)
	r.Get("/:userId", getUserHandler)
	r.Post("/", createUserHandler)
	r.Put("/", updateUserHandler)

	return r
}
