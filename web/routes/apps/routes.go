package apps

import (
	"net/http"

	"github.com/pressly/chi"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	// r.Get("/", getUsersHandler)
	// r.Get("/:userId", getUserHandler)
	// r.With(web.BodyParser(&createUser{}, 256)).Post("/", createUserHandler)
	// r.With(web.BodyParser(&updateUser{}, 256)).Put("/", updateUserHandler)

	return r
}
