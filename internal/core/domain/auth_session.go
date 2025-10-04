package domain

import "time"

type SessionID string
type State string // A string to protect against CSRF attacks
type Nonce string // A string used to protect against replay attacks

// AuthorizationSession stores the parameters of the original authorization request
// is created in the first step and is used for verification in subsequent steps
// ensures security before code is released
type AuthorizationSession struct {
	ID          SessionID `json:"id"`
	ClientID    ClientID  `json:"client_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       []Scope   `json:"scope"`
	State       State     `json:"state"`
	Nonce       *Nonce    `json:"nonce"` // Associates a client session with an ID Token to protect against replay attacks
	PKCE        *PKCE     `json:"pkce"`
	Provider    Provider  `json:"provider"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"` // The session should have a short lifetime (3-5 minutes)
}
