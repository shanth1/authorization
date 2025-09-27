package domain

import "time"

type TokenID string
type JTI string

type AccessToken struct {
	ID        TokenID   `json:"id"`
	UserID    UserID    `json:"user_id"`
	ClientID  ClientID  `json:"client_id"`
	Scope     []Scope   `json:"scope"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
	JTI       JTI       `json:"jti"`
}

type RefreshToken struct {
	ID          TokenID   `json:"id"`
	UserID      UserID    `json:"user_id"`
	ClientID    ClientID  `json:"client_id"`
	Scope       []Scope   `json:"scope"`
	ExpiresAt   time.Time `json:"expires_at"`
	IssuedAt    time.Time `json:"issued_at"`
	RotatedFrom *TokenID  `json:"rotated_from"` // reference to the old refresh token
	JTI         JTI       `json:"jti"`
}

type IDTokenClaims struct {
	Issuer        string    `json:"issuer"`
	Subject       string    `json:"subject"`
	Audience      string    `json:"audience"`
	Expires       time.Time `json:"expires"`
	IssuedAt      time.Time `json:"issued_at"`
	Nonce         *Nonce    `json:"nonce"`
	Email         *string   `json:"email"`
	EmailVerified *bool     `json:"email_verified"`
	Name          *string   `json:"name"`
	Picture       *string   `json:"picture"`
}

type LogoutTokenClaims struct {
	Issuer   string         `json:"issuer"`
	Subject  string         `json:"subject"`
	Audience string         `json:"audience"`
	IssuedAt time.Time      `json:"issued_at"`
	JTI      JTI            `json:"jti"`
	Events   map[string]any `json:"events"` // per OIDC logout event
	SID      *string        `json:"sid"`    // optional session id
}
