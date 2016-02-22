package apps

import (
	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web/middleware"
)

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.JwtHandler())

	r.Get("/", listAllAppsHandler)
	r.Post("/", middleware.BodyParser(appRequestBuilder, 512), createAppHandler)
	r.Patch("/:userId", middleware.BodyParser(appRequestBuilder, 512), updateAppHandler)

	return r
}
