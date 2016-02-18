package session

import (
	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web/middleware"
)

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/start", middleware.BodyParser(loginRequestBuilder, 512), startSessionHandler)
	r.Get("/end", endSessionHandler)

	return r
}
