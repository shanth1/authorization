package domain

import "time"

type SessionID string
type State string

// AuthorizationSession stores the parameters of the original authorization request.
// ensures security before code is released
type AuthorizationSession struct {
	ID          SessionID `json:"id"`
	ClientID    ClientID  `json:"client_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       []Scope   `json:"scope"`
	State       State     `json:"state"`
	Nonce       *Nonce    `json:"nonce"`
	PKCE        *PKCE     `json:"pkce"`
	Provider    Provider  `json:"provider"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
