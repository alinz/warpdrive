package util

import "net/http"

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
