package session

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.With(web.BodyParser(&sessionLogin{}, 256)).Post("/start", startSessionHandler)
	r.Get("/end", endSessionHandler)

	r.Group(func(r chi.Router) {
		r.Use(web.TokenAuth.Verifier)
		r.Use(web.Authenticator)

		r.Get("/", validateSessionHandler)
	})

	return r
}
