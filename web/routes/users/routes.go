package users

import (
	"net/http"

	"github.com/pressly/chi"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	// r.With(web.BodyParser(&sessionLogin{}, 256)).Post("/start", startSessionHandler)
	r.Get("/", getUserProfileHandler)

	return r
}
