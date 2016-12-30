package apps

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", getAppsHandler)
	r.Get("/:appId", getAppHandler)
	r.With(web.BodyParser(&createApp{}, 256)).Post("/", createAppHandler)
	r.With(web.BodyParser(&updateApp{}, 256)).Put("/:appId", updateAppHandler)

	r.Route("/:appId", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Get("/", getUsersAppHandler)
			r.Route("/:userId", func(r chi.Router) {
				r.Post("/", assignUserToAppHandler)
				r.Delete("/", unassignUserFromAppHandler)
			})
		})

		r.Route("/cycles", func(r chi.Router) {
			r.Get("/", getCyclesHandler)
			r.With(web.BodyParser(&createCycle{}, 128)).Post("/", createCycleHandler)

			r.Route("/:cycleId", func(r chi.Router) {
				r.Get("/", getCycleHandler)
				r.With(web.BodyParser(&updateCycle{}, 128)).Put("/", updateCycleHandler)
				r.Delete("/", removeCycleHandler)
				r.Get("/key", getCycleKeyHandler)

				r.Route("/releases", func(r chi.Router) {
					r.With(web.BodyParser(&createRelease{}, 1024)).Post("/", createReleaseHandler)
					r.Route("/:releaseId", func(r chi.Router) {
						r.Get("/", getReleaseHandler)
						r.With(web.BodyParser(&updateRelease{}, 1024)).Put("/", updateReleaseHandler)
						r.Delete("/", removeReleaseHandler)

						r.Route("/bundles", func(r chi.Router) {
							r.Post("/", uploadBundlesHandler)
							r.Get("/", getBundlesHandler)
						})

						r.Route("/lock", func(r chi.Router) {
							r.Post("/", lockReleaseHandler)
							r.Delete("/", unlockReleaseHandler)
						})
					})
				})
			})
		})

	})

	return r
}
