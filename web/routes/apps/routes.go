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

	r.Route("/:appId/users", func(r chi.Router) {
		r.Get("/", usersAppHandler)
		r.Route("/:userId", func(r chi.Router) {
			r.Post("/", assignUserToAppHandler)
			r.Delete("/", unassignUserFromAppHandler)
		})
	})

	return r
}
