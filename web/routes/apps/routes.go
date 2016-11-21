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
	// r.With(web.BodyParser(&updateUser{}, 256)).Put("/", updateUserHandler)

	return r
}
