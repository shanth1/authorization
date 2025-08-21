package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shanth1/authorization/internal/ports"
)

// English comment: AppBackend service. This acts as a confidential client in OIDC, handling token exchanges, refreshes, user info proxies, and logouts. It validates tokens locally using cached JWKS and proxies requests to AuthBackend with client_secret for security.

type AppService struct {
	clientID     string
	clientSecret string
	authURL      string // Base URL for AuthBackend
	cache        ports.CachePort
	httpClient   *http.Client           // For requests to AuthBackend
	jwksKeys     map[string]interface{} // Parsed JWKS keys, loaded in FetchJWKS
}

// NewAppService creates a new AppService instance.
func NewAppService(clientID, clientSecret, authURL string, cache ports.CachePort) *AppService {
	return &AppService{
		clientID:     clientID,
		clientSecret: clientSecret,
		authURL:      authURL,
		cache:        cache,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		jwksKeys:     make(map[string]interface{}),
	}
}

// FetchJWKS fetches and caches JWKS from AuthBackend for local token validation.
func (s *AppService) FetchJWKS() error {
	resp, err := s.httpClient.Get(s.authURL + "/.well-known/jwks.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	for _, key := range jwks.Keys {
		kid, _ := key["kid"].(string)
		s.jwksKeys[kid] = key // Store full key map; in prod, parse to rsa.PublicKey
	}
	// Cache the JWKS with TTL (e.g., 1h)
	s.cache.Set("jwks", jwks, 1*time.Hour)
	return nil
}

// validateToken locally validates a JWT access_token using cached JWKS.
func (s *AppService) validateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // For demo HS256; change to RS256 in prod
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// For HS256 demo, use shared secret; in prod, use public key from JWKS
		return []byte("supersecretkey"), nil // Match AuthBackend privateKey
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !claims.VerifyAudience(s.clientID, true) || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("invalid claims")
	}
	return token, nil
}

// ExchangeCodeForToken handler: Exchanges auth_code for tokens via AuthBackend.
func (s *AppService) ExchangeCodeForToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	codeVerifier := r.FormValue("code_verifier")

	// Prepare request to AuthBackend /token
	data := strings.NewReader(fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s&code_verifier=%s",
		code, "http://localhost:8081/", s.clientID, s.clientSecret, codeVerifier)) // redirect_uri hardcoded for demo
	req, err := http.NewRequest("POST", s.authURL+"/token", data)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "token exchange failed", http.StatusUnauthorized)
		return
	}
	defer resp.Body.Close()

	var tokens map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		http.Error(w, "invalid response", http.StatusInternalServerError)
		return
	}

	// Set secure cookies (HttpOnly for security)
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: tokens["access_token"].(string), HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "id_token", Value: tokens["id_token"].(string), HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: tokens["refresh_token"].(string), HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, Path: "/"})

	// Return tokens to frontend (for demo; in prod, rely on cookies)
	json.NewEncoder(w).Encode(tokens)
}

// RefreshToken handler: Refreshes access_token using refresh_token, with rotation.
func (s *AppService) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	refreshToken := r.FormValue("refresh_token")

	// Prepare request to AuthBackend /token with grant_type=refresh_token
	data := strings.NewReader(fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s",
		refreshToken, s.clientID, s.clientSecret))
	req, err := http.NewRequest("POST", s.authURL+"/token", data)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "refresh failed", http.StatusUnauthorized)
		return
	}
	defer resp.Body.Close()

	var tokens map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		http.Error(w, "invalid response", http.StatusInternalServerError)
		return
	}

	// Update secure cookies with new tokens (rotation: new refresh_token)
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: tokens["access_token"].(string), HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: tokens["refresh_token"].(string), HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, Path: "/"})

	json.NewEncoder(w).Encode(tokens)
}

// GetUserInfo handler: Proxies to AuthBackend /userinfo, validates access_token locally first.
func (s *AppService) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Local validation
	_, err := s.validateToken(tokenStr)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Proxy to AuthBackend
	req, err := http.NewRequest("GET", s.authURL+"/userinfo", nil)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "userinfo failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response to client
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// Logout handler: Initiates global logout via AuthBackend, clears cookies.
func (s *AppService) Logout(w http.ResponseWriter, r *http.Request) {
	// Get access_token from cookie or header (for demo, from cookie)
	accessCookie, _ := r.Cookie("access_token")
	if accessCookie == nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	// Send to AuthBackend /logout
	data := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s", s.clientID, s.clientSecret))
	req, err := http.NewRequest("POST", s.authURL+"/logout", data)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessCookie.Value) // Pass access_token for user_id extraction

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}
	resp.Body.Close()

	// Clear cookies
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", MaxAge: -1, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "id_token", Value: "", MaxAge: -1, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", MaxAge: -1, Path: "/"})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logged out"))
}
