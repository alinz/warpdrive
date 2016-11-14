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

	return r
}
