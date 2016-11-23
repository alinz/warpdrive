package session

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goware/jwtauth"
	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
	"golang.org/x/crypto/bcrypt"
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err.Error())
		web.Respond(w, http.StatusUnauthorized, nil)
		return
	}

	var claims jwtauth.Claims
	claims = make(map[string]interface{})
	claims.Set("user:id", fmt.Sprintf("%v", user.ID))

	_, token, _ := web.TokenAuth.Encode(claims)

	web.SetJWTCookie(w, r, token)
	web.Respond(w, 200, nil)
}

func endSessionHandler(w http.ResponseWriter, r *http.Request) {
	web.SetJWTCookie(w, r, "")
	web.Respond(w, 200, nil)
}

func validateSessionHandler(w http.ResponseWriter, r *http.Request) {
	web.Respond(w, 200, nil)
}
