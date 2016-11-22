package routes

import "net/http"

func checkVersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}
