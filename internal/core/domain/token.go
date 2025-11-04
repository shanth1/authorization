package domain

import "time"

type TokenID string
type JTI string
type URL string

// AccessToken used for requests to protected resources
type AccessToken struct {
	ID        TokenID   `json:"id"`
	UserID    UserID    `json:"user_id"`
	ClientID  ClientID  `json:"client_id"`
	Scope     []Scope   `json:"scope"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
	JTI       JTI       `json:"jti"`
}

// RefreshToken used to obtain a new access/refresh token pair without re-authenticating the user
type RefreshToken struct {
	ID            TokenID       `json:"id"`
	UserID        UserID        `json:"user_id"`
	ClientID      ClientID      `json:"client_id"`
	UserSessionID UserSessionID `json:"user_session_id"`
	Scope         []Scope       `json:"scope"`
	ExpiresAt     time.Time     `json:"expires_at"`
	IssuedAt      time.Time     `json:"issued_at"`
	RotatedFrom   *TokenID      `json:"rotated_from"` // reference to the old refresh token
	JTI           JTI           `json:"jti"`
	Revoked       bool          `json:"revoked"`
	RevokedAt     *time.Time    `json:"revoked_at"`
}

// IDTokenClaims is a JWT that verifies the user's identity for the client backend
type IDTokenClaims struct {
	Issuer        URL        `json:"iss"` // Authorization server URL
	Subject       UserID     `json:"sub"`
	Audience      []ClientID `json:"aud"`
	Expires       time.Time  `json:"exp"`
	IssuedAt      time.Time  `json:"iat"`
	AuthTime      time.Time  `json:"auth_time"`
	Nonce         *Nonce     `json:"nonce"`
	Email         *string    `json:"email"`
	EmailVerified *bool      `json:"email_verified"`
	Name          *string    `json:"name"`
	Picture       *string    `json:"picture"`
}

type LogoutTokenClaims struct {
	Issuer    URL            `json:"iss"` // Authorization server URL
	Audience  string         `json:"aud"`
	IssuedAt  time.Time      `json:"iat"`
	JTI       JTI            `json:"jti"`           // Unique JWT identifier to prevent reuse
	Events    map[string]any `json:"events"`        // per OIDC logout event
	Subject   *UserID        `json:"sub,omitempty"` // [optional] required if session id is not specified
	SessionID *UserSessionID `json:"sid,omitempty"` // [optional] required if subject is not specified
}
