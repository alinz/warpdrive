package constant

import (
	"errors"
	"net/http"
)

var (
	ErrUnauthorized      = errors.New("unauthorized token")
	ErrorAuthorizeAccess = errors.New("unauthorized access")
)

func ErrorStatusCode(err error) int {
	switch err {
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrorAuthorizeAccess:
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}
