// internal/app/auth_service.go
package app

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/shanth1/authorization/internal/domain"
	"github.com/shanth1/authorization/internal/ports"
)

// English comment: Application service for AuthBackend, hexagonal core logic.

type AuthService struct {
	userRepo     ports.UserRepository
	clientRepo   ports.ClientRepository
	sessionRepo  ports.SessionRepository
	tokenRepo    ports.TokenRepository
	cache        ports.CachePort
	providerPort ports.ProviderPort
	jwks         *JWKS // Struct with private key, etc.
}

func NewAuthService(userRepo ports.UserRepository, clientRepo ports.ClientRepository, sessionRepo ports.SessionRepository, tokenRepo ports.TokenRepository, cache ports.CachePort, providerPort ports.ProviderPort, jwks *JWKS) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		clientRepo:   clientRepo,
		sessionRepo:  sessionRepo,
		tokenRepo:    tokenRepo,
		cache:        cache,
		providerPort: providerPort,
		jwks:         jwks,
	}
}

// RegisterClient handler
func (s *AuthService) RegisterClient(w http.ResponseWriter, r *http.Request) {
	// Parse body: name, redirect_uris, etc.
	var req struct {
		Name         string   `json:"name"`
		RedirectURIs []string `json:"redirect_uris"`
		Type         string   `json:"type"`
		Scopes       []string `json:"scopes"`
		LogoutURI    string   `json:"logout_uri"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	client := &domain.Client{
		ID:           uuid.New(),
		Name:         req.Name,
		Secret:       generateRandomString(32), // If confidential
		RedirectURIs: req.RedirectURIs,
		Type:         req.Type,
		Scopes:       req.Scopes,
		LogoutURI:    req.LogoutURI,
	}
	s.clientRepo.Create(client)

	json.NewEncoder(w).Encode(map[string]string{"client_id": client.ID.String(), "client_secret": client.Secret})
}

// GetProviders
func (s *AuthService) GetProviders(w http.ResponseWriter, r *http.Request) {
	// For scalability, return list from config or DB
	providers := []string{"login_password"} // Add more in future
	json.NewEncoder(w).Encode(providers)
}

// Authorize handler
func (s *AuthService) Authorize(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Render form, but since JS, assume POST for auth
		return
	}

	// Parse form
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	state := r.FormValue("state")
	codeChallenge := r.FormValue("code_challenge")
	nonce := r.FormValue("nonce")
	// provider := r.FormValue("provider")
	login := r.FormValue("login")
	password := r.FormValue("password")

	// Validate client
	client, err := s.clientRepo.FindByID(clientID)
	if err != nil || !contains(client.RedirectURIs, redirectURI) {
		http.Error(w, "invalid client", http.StatusBadRequest)
		return
	}

	// Authenticate via provider
	identity, err := s.providerPort.Authenticate(map[string]string{"login": login, "password": password})
	if err != nil {
		http.Error(w, "auth failed", http.StatusUnauthorized)
		return
	}

	// Find or create user, merge if needed
	user, err := s.userRepo.FindByEmail(identity.Claims["email"].(string))
	if err != nil {
		user = &domain.User{
			ID:        uuid.New(),
			Email:     identity.Claims["email"].(string),
			Providers: []domain.ProviderIdentity{*identity},
		}
		s.userRepo.Create(user)
	} else {
		// Merge if not exists
		exists := false
		for _, p := range user.Providers {
			if p.ProviderType == identity.ProviderType {
				exists = true
				break
			}
		}
		if !exists {
			user.Providers = append(user.Providers, *identity)
			s.userRepo.Update(user)
		}
	}

	// Create session
	sessionID := uuid.New()
	session := &domain.Session{
		ID:            sessionID,
		UserID:        user.ID,
		ClientID:      client.ID,
		State:         state,
		Nonce:         nonce,
		CodeChallenge: codeChallenge,
		ExpiresAt:     time.Now().Add(10 * time.Minute),
	}
	s.sessionRepo.Create(session)

	// Generate code
	code := generateRandomString(32)
	token := &domain.Token{
		Value:     code,
		Type:      "code",
		UserID:    user.ID,
		ClientID:  client.ID,
		Scope:     strings.Join(client.Scopes, " "),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	s.tokenRepo.Create(token)

	// Redirect
	redirect := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)
	http.Redirect(w, r, redirect, http.StatusFound)
}

// Token handler
func (s *AuthService) Token(w http.ResponseWriter, r *http.Request) {
	grantType := r.FormValue("grant_type")
	if grantType == "authorization_code" {
		code := r.FormValue("code")
		codeVerifier := r.FormValue("code_verifier")
		clientID := r.FormValue("client_id")
		clientSecret := r.FormValue("client_secret")
		redirectURI := r.FormValue("redirect_uri")

		// Validate
		token, err := s.tokenRepo.FindByValue(code)
		if err != nil || token.Type != "code" || time.Now().After(token.ExpiresAt) {
			http.Error(w, "invalid code", http.StatusBadRequest)
			return
		}
		client, _ := s.clientRepo.FindByID(clientID)
		if client.Secret != clientSecret || !contains(client.RedirectURIs, redirectURI) {
			http.Error(w, "invalid client", http.StatusBadRequest)
			return
		}
		session, _ := s.sessionRepo.FindByID(token.UserID.String()) // Assume session linked
		challenge := sha256.Sum256([]byte(codeVerifier))
		if base64.URLEncoding.EncodeToString(challenge[:]) != session.CodeChallenge {
			http.Error(w, "invalid pkce", http.StatusBadRequest)
			return
		}

		// Generate tokens
		accessToken, _ := s.generateJWT(token.UserID.String(), clientID, token.Scope, session.Nonce, 1*time.Hour)
		idToken, _ := s.generateIDToken(token.UserID.String(), clientID, session.Nonce, time.Now().Unix())
		refreshToken := uuid.New().String()
		refresh := &domain.Token{
			Value:     refreshToken,
			Type:      "refresh",
			UserID:    token.UserID,
			ClientID:  token.ClientID,
			Scope:     token.Scope,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}
		s.tokenRepo.Create(refresh)

		// Delete code
		s.tokenRepo.DeleteByValue(code)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  accessToken,
			"id_token":      idToken,
			"refresh_token": refreshToken,
			"expires_in":    3600,
		})
	} else if grantType == "refresh_token" {
		// Similar logic: validate, rotate refresh, new access
		// ...
	}
}

// Other methods: UserInfo, Logout, JWKS similarly implemented

// Helpers
func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// JWKS struct and generateJWT (use jwt lib)
type JWKS struct {
	privateKey []byte
}

func NewJWKS(private []byte) *JWKS {
	return &JWKS{privateKey: private}
}

func (s *AuthService) generateJWT(sub, aud, scope, nonce string, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   sub,
		"aud":   aud,
		"scope": scope,
		"exp":   time.Now().Add(exp).Unix(),
		"iat":   time.Now().Unix(),
		"nonce": nonce,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // HS for demo, RS in prod
	return token.SignedString(s.jwks.privateKey)
}

// Similar for idToken

func (s *AuthService) JWKS(w http.ResponseWriter, r *http.Request) {
	// Return public keys
	json.NewEncoder(w).Encode(map[string]interface{}{"keys": []interface{}{ /* public key */ }})
}

// Logout
func (s *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	// Parse token, get userID
	// s.sessionRepo.DeleteByUserID(userID)
	// s.tokenRepo.DeleteByUserID(userID)
	// Notify clients via HTTP POST to logoutURI
	// Redirect to AuthFrontend
}

func (s *AuthService) UserInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate access_token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwks.privateKey, nil // HS256 for demo
	})
	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "invalid claims", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		http.Error(w, "invalid sub", http.StatusUnauthorized)
		return
	}

	// Find user
	user, err := s.userRepo.FindByEmail(userID.String()) // Add FindByID to UserRepository if not
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Return claims
	info := map[string]interface{}{
		"sub":       user.ID.String(),
		"email":     user.Email,
		"providers": user.Providers, // For merge info
		// Add more profile claims if scope allows
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	json.NewEncoder(w).Encode(info)
}

func (s *AuthService) generateIDToken(sub, aud, nonce string, authTime int64) (string, error) {
	claims := jwt.MapClaims{
		"iss":       "http://localhost:8082", // Issuer
		"sub":       sub,
		"aud":       aud,
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
		"nonce":     nonce,
		"auth_time": authTime,
		// at_hash, acr, etc. if needed
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwks.privateKey)
}
