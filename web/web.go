package web

import (
	"net/http"

	"github.com/pressly/chi"
)

func New() http.Handler {
	r := chi.NewRouter()

	return r
}
