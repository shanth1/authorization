package server

import "net/http"

func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	default:
		http.Error(w, `{"error": "Not found"}`, http.StatusNotFound)
	}
}
