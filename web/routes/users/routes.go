package users

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(web.TokenAuth.Verifier)
	r.Use(web.Authenticator)

	r.Get("/", getUsersHandler)
	r.Get("/:userId", getUserHandler)
	r.With(web.BodyParser(&createUser{}, 256)).Post("/", createUserHandler)
	r.With(web.BodyParser(&updateUser{}, 256)).Put("/", updateUserHandler)

	return r
}
