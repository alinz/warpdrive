package web

import (
	"net/http"

	"github.com/goware/heartbeat"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"

	"github.com/pressly/warpdrive/web/handler/apps"
	"github.com/pressly/warpdrive/web/handler/session"
	"github.com/pressly/warpdrive/web/handler/users"
)

func New() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Logger)

	r.Use(heartbeat.Route("/ping"))

	r.Mount("/session", session.Routes())
	r.Mount("/users", users.Routes())
	r.Mount("/apps", apps.Routes())

	return r
}
