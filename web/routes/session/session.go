package session

import (
	"net/http"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

type sessionLogin struct {
	Email    *string `json:"email,required"`
	Password *string `json:"password,required"`
}

func startSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := ctx.Value("parsed:body").(*sessionLogin)
	email := *body.Email
	password := *body.Password

	if err := services.FindUserByEmailPassword(email, password); err != nil {
		statusCode := statusCodeError(err)
		web.Respond(w, statusCode, err)
		return
	}

	web.Respond(w, 200, nil)
}

func endSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`end`))
}
