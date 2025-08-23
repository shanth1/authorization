package server

import (
	"net/http"
	"strings"

	"github.com/shanth1/authorization/client/internal/config"
)

type handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *handler {
	return &handler{
		cfg: cfg,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, "/api") {
		h.handleAPI(w, r)
		return
	}

	h.serveStatic(w, r)
}
