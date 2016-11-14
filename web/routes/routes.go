package routes

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/web"
	"github.com/pressly/warpdrive/web/routes/session"
	"github.com/pressly/warpdrive/web/routes/users"
)

func New() chi.Router {
	conf := warpdrive.Conf
	web.JwtSetup(conf.JWT.SecretKey)

	r := chi.NewRouter()

	r.Get("/", index)

	r.Mount("/session", session.Routes())

	r.Group(func(r chi.Router) {
		r.Use(web.TokenAuth.Verifier)
		r.Use(web.Authenticator)

		r.Mount("/users", users.Routes())
	})

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
