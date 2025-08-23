package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	baseDir := "./client/static"
	fullPath := filepath.Join(baseDir, filepath.Clean(path))

	baseDirBase := filepath.Base(baseDir)
	fullPathBase := filepath.Base(filepath.Dir(fullPath))

	if !strings.HasPrefix(fullPathBase, baseDirBase) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, fullPath)
}
