package middleware

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web/util"
	"golang.org/x/net/context"
)

//BodyParser loads builder with maxSize and tries to load the message.
//if for some reason it can't parse the message, it will return an error.
//if successful, it will put the processed data into context with key 'json_body'
func BodyParser(builder func() interface{}, maxSize int64) func(chi.Handler) chi.Handler {
	return func(next chi.Handler) chi.Handler {
		return chi.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			to := builder()

			if err := util.StreamJSONToStructWithLimit(r.Body, to, maxSize); err != nil {
				util.Respond(w, 422, err)
				return
			}

			//check for required fields
			if err := util.JSONValidation(to); err != nil {
				util.Respond(w, 400, err)
				return
			}

			ctx = context.WithValue(ctx, constant.CtxKeyParsedBody, to)

			next.ServeHTTPC(ctx, w, r)
		})
	}
}
