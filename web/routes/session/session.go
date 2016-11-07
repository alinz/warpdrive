package session

import (
	"net/http"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func startSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`start`))

	email := ""
	password := ""

	if err := services.CreateSession(email, password); err != nil {
		statusCode := statusCodeError(err)
		web.Respond(w, statusCode, err)
	}
}

func endSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`end`))

	services.DestorySession()
}
