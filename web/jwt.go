package web

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/goware/jwtauth"
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

		userID, ok := jwtToken.Claims["userid"]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "userid", userID.(int))

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
