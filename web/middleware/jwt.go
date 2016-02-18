package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/pressly/warpdrive/web/security"

	"github.com/pressly/chi"
	"golang.org/x/net/context"
)

//JwtHandler this is middleware for chi to inject jwt token into context
func JwtHandler() func(chi.Handler) chi.Handler {
	return func(next chi.Handler) chi.Handler {
		hfn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

			token, err := security.TryFindJwt(r)

			if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(401)
				json.NewEncoder(w).Encode(err)
				return
			}

			ctx = context.WithValue(ctx, "jwt", token.Raw)
			ctx = context.WithValue(ctx, "jwt.token", token)

			next.ServeHTTPC(ctx, w, r)
		}
		return chi.HandlerFunc(hfn)
	}
}
