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
