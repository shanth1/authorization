package server

import (
	"encoding/json"
	"net/http"
)

func (h *handler) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/api/config":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		config := map[string]string{
			"authUrl":  h.cfg.AuthURL,
			"clientId": h.cfg.ClientID,
		}

		json.NewEncoder(w).Encode(config)
	default:
		http.Error(w, `{"error": "Not found"}`, http.StatusNotFound)
	}
}
