package routes

import (
	"net/http"

	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/web/routes/session"
)

func New() chi.Router {
	conf := warpdrive.Conf
	tokenAuth := jwtauth.New("HS256", []byte(conf.JWT.SecretKey), nil)

	r := chi.NewRouter()

	r.Get("/", index)

	r.Mount("/session", session.Routes())

	r.Group(func(r chi.Router) {
		r.Use(tokenAuth.Verifier)

	})

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
