// cmd/auth_backend/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanth1/authorization/internal/adapters"
	"github.com/shanth1/authorization/internal/app"
)

// English comment: Main entry for AuthBackend, sets up HTTP server with routes.
func main() {
	// Initialize repositories (ports implementations)
	userRepo := adapters.NewInMemoryUserRepository()
	clientRepo := adapters.NewInMemoryClientRepository()
	sessionRepo := adapters.NewInMemorySessionRepository()
	tokenRepo := adapters.NewInMemoryTokenRepository()
	cache := adapters.NewInMemoryCache()
	providerPort := adapters.NewLoginPasswordProvider(userRepo) // Only login/password for example

	// JWKS keys (generate RSA key pair in production)
	privateKey := []byte("supersecretkey") // For demo, use real RSA in prod
	jwks := app.NewJWKS(privateKey)

	// Application service
	authService := app.NewAuthService(userRepo, clientRepo, sessionRepo, tokenRepo, cache, providerPort, jwks)

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/register-client", authService.RegisterClient).Methods("POST")
	r.HandleFunc("/providers", authService.GetProviders).Methods("GET")
	r.HandleFunc("/authorize", authService.Authorize).Methods("GET", "POST")
	r.HandleFunc("/token", authService.Token).Methods("POST")
	r.HandleFunc("/userinfo", authService.UserInfo).Methods("GET")
	r.HandleFunc("/logout", authService.Logout).Methods("POST")
	r.HandleFunc("/.well-known/jwks.json", authService.JWKS).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/auth_frontend"))) // Serve AuthFrontend

	log.Println("AuthBackend starting on :8082")
	http.ListenAndServe(":8082", r)
}
