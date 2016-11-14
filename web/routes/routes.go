package routes

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/web/routes/session"
)

func authenticator(next http.Handler) http.Handler {
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

func New() chi.Router {
	conf := warpdrive.Conf
	tokenAuth := jwtauth.New("HS256", []byte(conf.JWT.SecretKey), nil)

	r := chi.NewRouter()

	r.Get("/", index)

	r.Mount("/session", session.Routes())

	r.Group(func(r chi.Router) {
		r.Use(tokenAuth.Verifier)
		r.Use(authenticator)
	})

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
