package session

import (
	"net/http"

	"github.com/pressly/chi"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/start", startSessionHandler)
	r.Get("/end", endSessionHandler)

	return r
}
