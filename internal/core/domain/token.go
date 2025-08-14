package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Tokens - набор токенов OAuth 2.0 / OIDC
type Tokens struct {
	AccessToken  string `json:"access_token"`  // Токен доступа
	RefreshToken string `json:"refresh_token"` // Токен для обновления
	IDToken      string `json:"id_token"`      // JWT с информацией о пользователе
	ExpiresIn    int64  `json:"expires_in"`    // Время жизни в секундах
}

// TokenClaims - кастомные JWT-claims для IdP
type TokenClaims struct {
	jwt.RegisteredClaims
	UserID   string   `json:"sub"`    // Subject (идентификатор пользователя)
	ClientID string   `json:"aud"`    // Audience (идентификатор клиента)
	Nonce    string   `json:"nonce"`  // Уникальный идентификатор сессии
	Scopes   []string `json:"scopes"` // Разрешения
}

// Valid реализует интерфейс jwt.Claims
func (c *TokenClaims) Valid() error {
	if c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now()) {
		return jwt.ErrTokenExpired
	}
	return nil
}

// GetAudience возвращает аудиторию токена
func (c TokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{c.ClientID}, nil
}

// GetSubject возвращает субъект токена (реализация SubjectClaim)
func (c *TokenClaims) GetSubject() (string, error) {
	return c.UserID, nil
}

type TokenRequest struct {
	GrantType   string `form:"grant_type" json:"grant_type"`
	Code        string `form:"code" json:"code"`
	RedirectURI string `form:"redirect_uri" json:"redirect_uri"`
	ClientID    string `form:"client_id" json:"client_id"`

	// Для PKCE (опционально)
	CodeVerifier string `form:"code_verifier" json:"code_verifier"`
}
