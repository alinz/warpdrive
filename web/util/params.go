package util

import (
	"strconv"

	"github.com/pressly/chi"
	"golang.org/x/net/context"
)

//ParamValue accepts the context and returns param's value
func ParamValue(ctx context.Context, param string) string {
	value := chi.URLParams(ctx)[param]
	return value
}

//ParamValueAsID accepts the context and try to parse the parama into int64 id
func ParamValueAsID(ctx context.Context, param string) (int64, error) {
	id, err := strconv.ParseInt(ParamValue(ctx, param), 10, 64)

	if err != nil {
		return -1, err
	}

	return id, nil
}
