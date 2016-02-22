package web

import (
	"net/http"

	"github.com/pressly/warpdrive/web/handler/apps"
	"github.com/pressly/warpdrive/web/handler/session"
	"github.com/pressly/warpdrive/web/handler/users"

	"github.com/pressly/chi"
)

func New() http.Handler {
	r := chi.NewRouter()

	r.Mount("/session", session.Routes())
	r.Mount("/users", users.Routes())
	r.Mount("/apps", apps.Routes())

	return r
}
