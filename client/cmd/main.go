package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/", handleRequest)

	addr := ":8100"
	fmt.Printf("Server started at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, "/api") {
		handleAPI(w, r)
		return
	}

	serveStatic(w, r)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	default:
		http.Error(w, `{"error": "Not found"}`, http.StatusNotFound)
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	baseDir := "./static"
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
