package util

import (
	"net/http"

	"github.com/pressly/warpdrive/web/constant"
)

//Respond is a utility function which helps identifies error from other messages
func Respond(w http.ResponseWriter, status int, v interface{}) {
	if err, ok := v.(error); ok {
		message := map[string]interface{}{"error": err.Error()}
		WriteAsJSON(w, message, status)
		return
	}

	if v != nil {
		WriteAsJSON(w, v, status)
	} else {
		w.WriteHeader(status)
	}
}

//RespondError is helper function to return proper status code with error message
func RespondError(w http.ResponseWriter, err error) {
	status := constant.ErrorStatusCode(err)

	if err != nil {
		message := map[string]interface{}{"error": err.Error()}
		WriteAsJSON(w, message, status)
	} else {
		WriteAsJSON(w, nil, status)
	}
}
