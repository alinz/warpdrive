package web

import (
	"context"
	"net/http"
	"strings"
	"time"

	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/goware/jwtauth"
	"github.com/pressly/warpdrive"
)

var (
	TokenAuth *jwtauth.JwtAuth
)

func JwtSetup(secretKey string) {
	if TokenAuth != nil {
		panic("Jwt Called twice!")
	}

	TokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if jwtErr, ok := ctx.Value("jwt.err").(error); ok {
			if jwtErr != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}

		jwtToken, ok := ctx.Value("jwt").(*jwt.Token)
		if !ok || jwtToken == nil || !jwtToken.Valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userID, ok := jwtToken.Claims["userId"]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		id, _ := strconv.ParseInt(userID.(string), 10, 64)
		ctx = context.WithValue(ctx, "userId", id)

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetJWTCookie(w http.ResponseWriter, r *http.Request, value string) {
	path := warpdrive.Conf.JWT.Path
	maxAge := warpdrive.Conf.JWT.MaxAge
	expires := time.Now().AddDate(1, 0, 0) // expiry date in 1 year

	if value == "" {
		maxAge = -1               // delete cookie now
		expires = time.Unix(1, 0) // set to epoche for delete
	}

	host := warpdrive.Conf.JWT.Domain
	if strings.Index(host, "localhost") >= 0 {
		host = ""
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Domain:   host,
		Path:     path,
		Secure:   warpdrive.Conf.JWT.Secure,
		MaxAge:   maxAge,
		Expires:  expires,
		HttpOnly: true,
		Value:    value,
	}
	http.SetCookie(w, &cookie)
}
