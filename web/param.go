package web

import (
	"net/http"
	"strconv"

	"github.com/pressly/chi"
)

func ParamAsInt64(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, name), 10, 64)
}

func ParamAsString(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}
