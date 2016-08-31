package session

import "net/http"

func startSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`start`))
}

func endSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`end`))
}
