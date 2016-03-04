package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pressly/chi"
	"golang.org/x/net/context"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//GZipHandler supports gzip feature for chi
func GZipHandler() func(chi.Handler) chi.Handler {
	return func(next chi.Handler) chi.Handler {
		hfn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Header.Get("Accept-Encoding"))
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				gz := gzip.NewWriter(w)
				defer gz.Close()

				w.Header().Set("Content-Encoding", "gzip")
				w = gzipResponseWriter{Writer: gz, ResponseWriter: w}
			}

			next.ServeHTTPC(ctx, w, r)
		}
		return chi.HandlerFunc(hfn)
	}
}
