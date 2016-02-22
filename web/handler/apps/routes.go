package apps

import (
	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web/middleware"
)

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.JwtHandler())

	r.Get("/", listAllAppsHandler)
	r.Post("/", middleware.BodyParser(createAppRequestBuilder, 512), createAppHandler)

	r.Route("/:appId", func(r chi.Router) {
		r.Patch("/", middleware.BodyParser(createAppRequestBuilder, 512), updateAppHandler)

		r.Post("/cycles", middleware.BodyParser(createCycleRequestBuilder, 512), createAppCycleHandler)
	})

	return r
}
