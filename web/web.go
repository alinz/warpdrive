package web

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"

	"github.com/pressly/warpdrive/web/handler/apps"
	"github.com/pressly/warpdrive/web/handler/session"
	"github.com/pressly/warpdrive/web/handler/users"
)

func New() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/session", session.Routes())
	r.Mount("/users", users.Routes())
	r.Mount("/apps", apps.Routes())

	return r
}
