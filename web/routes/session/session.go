package session

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

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

	user := services.FindUserByEmail(email)

	if user == nil {
		web.Respond(w, http.StatusUnauthorized, nil)
		return
	}

	fmt.Println(password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err.Error())
		web.Respond(w, http.StatusUnauthorized, nil)
		return
	}

	web.Respond(w, 200, nil)
}

func endSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`end`))
}
