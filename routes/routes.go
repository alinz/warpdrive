package routes

import (
	"net/http"

	"github.com/pressly/chi"
)

func New() chi.Router {
	r := chi.NewRouter()

	r.Get("/", index)

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
