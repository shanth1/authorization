package server

import (
	"net/http"
	"strings"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, "/api") {
		handleAPI(w, r)
		return
	}

	serveStatic(w, r)
}
