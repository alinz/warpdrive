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

// New this is the root level that creates all the routes for
// warpdrive server
func New(forDoc bool) chi.Router {
	if !forDoc {
		conf := warpdrive.Conf
		web.JwtSetup(conf.JWT.SecretKey)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Public routes
	r.Get("/", index)

	r.Mount("/session", session.Routes())
	r.Mount("/users", users.Routes())
	r.Mount("/apps", apps.Routes())

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
