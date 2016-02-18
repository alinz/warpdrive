package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web/constant"
	"github.com/pressly/warpdrive/web/util"
	"golang.org/x/net/context"
)

func ShouldBeRootHandler(next chi.Handler) chi.Handler {
	hfn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if !util.UserIsRoot(ctx) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(403)
			json.NewEncoder(w).Encode(constant.ErrorAuthorizeAccess)
			return
		}

		next.ServeHTTPC(ctx, w, r)
	}
	return chi.HandlerFunc(hfn)
}
