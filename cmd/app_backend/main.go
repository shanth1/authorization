// cmd/app_backend/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanth1/authorization/internal/adapters"
	"github.com/shanth1/authorization/internal/app"
)

// English comment: Main entry for AppBackend, example client backend.
func main() {
	// For demo, hardcoded client_id and secret from registration
	clientID := "example_client_id"
	clientSecret := "example_client_secret"
	authBackendURL := "http://localhost:8082"

	// Cache for JWKS
	cache := adapters.NewInMemoryCache()
	appService := app.NewAppService(clientID, clientSecret, authBackendURL, cache)

	r := mux.NewRouter()
	r.HandleFunc("/token", appService.ExchangeCodeForToken).Methods("POST")
	r.HandleFunc("/refresh", appService.RefreshToken).Methods("POST")
	r.HandleFunc("/userinfo", appService.GetUserInfo).Methods("GET")
	r.HandleFunc("/logout", appService.Logout).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/app_frontend"))) // Serve AppFrontend

	// Fetch and cache JWKS on start
	appService.FetchJWKS()

	log.Println("AppBackend starting on :8081")
	http.ListenAndServe(":8081", r)
}
