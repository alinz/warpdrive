package apps

import (
	"github.com/pressly/chi"
	m "github.com/pressly/warpdrive/web/middleware"
)

//Routes routes for app
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(m.JwtHandler())

		r.Get("/", listAllAppsHandler)
		r.Post("/", m.BodyParser(createAppRequestBuilder, 512), createAppHandler)

		r.Route("/:appId", func(r chi.Router) {
			r.Patch("/", m.BodyParser(createAppRequestBuilder, 512), updateAppHandler)

			r.Route("/cycles", func(r chi.Router) {
				r.Post("/", m.BodyParser(createCycleRequestBuilder, 512), createAppCycleHandler)
				r.Get("/", allAppCyclesHandler)

				r.Route("/:cycleId", func(r chi.Router) {
					r.Get("/", appCycle)
					r.Patch("/", m.BodyParser(updateCycleRequestBuilder, 512), updateAppCycleHandler)
					r.Get("/config", downloadAppCycleConfigHandler)

					r.Route("/releases", func(r chi.Router) {
						r.Post("/", createAppCycleReleaseHandler)
						r.Get("/", allAppCycleReleaseHandler)
						r.Patch("/:releaseId/lock", lockAppCycleReleaseHandler)
					})
				})
			})
		})
	})

	r.Group(func(r chi.Router) {
		r.Get("/:appId/cycles/:cycleId/releases/check", checkVersionAppCycleReleaseHandler)
		r.Post("/:appId/cycles/:cycleId/releases/download", m.GZipHandler(), downloadAppCycleReleaseHandler)
	})

	return r
}
