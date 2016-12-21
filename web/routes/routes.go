package routes

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/web"
	"github.com/pressly/warpdrive/web/routes/apps"
	"github.com/pressly/warpdrive/web/routes/session"
	"github.com/pressly/warpdrive/web/routes/users"
)

func New() chi.Router {
	conf := warpdrive.Conf
	web.JwtSetup(conf.JWT.SecretKey)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Public routes
	r.Get("/", index)

	r.Mount("/session", session.Routes())

	r.Get("/apps/:appId/cycles/:cycleId/releases/latest/version/:version/platform/:platform", checkVersionHandler)
	r.Post("/apps/:appId/cycles/:cycleId/releases/:releaseId/download", downloadHandler)

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(web.TokenAuth.Verifier)
		r.Use(web.Authenticator)

		r.Mount("/users", users.Routes())
		r.Mount("/apps", apps.Routes())
	})

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
