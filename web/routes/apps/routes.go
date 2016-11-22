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
			r.Get("/", usersAppHandler)
			r.Route("/:userId", func(r chi.Router) {
				r.Post("/", assignUserToAppHandler)
				r.Delete("/", unassignUserFromAppHandler)
			})
		})

		r.Route("/cycles", func(r chi.Router) {
			r.Get("/", cyclesAppHandler)
			r.With(web.BodyParser(&createCycle{}, 128)).Post("/", createCycleAppHandler)

			r.Route("/:cycleId", func(r chi.Router) {
				r.Get("/", getCycleAppHandler)
				r.Put("/", updateCycleAppHandler)
				r.Delete("/", removeCycleAppHandler)
				r.Get("/key", getKeyCycleAppHandler)

				r.Route("/releases", func(r chi.Router) {

				})
			})
		})
	})

	return r
}
