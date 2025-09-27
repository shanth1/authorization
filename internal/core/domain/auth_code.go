package domain

import "time"

type Nonce string
type CodeChallengeMethod string

const (
	CodeChallengeS256 CodeChallengeMethod = "S256"
)

type PKCE struct {
	CodeChallenge       string              `json:"code_challenge"`
	CodeChallengeMethod CodeChallengeMethod `json:"code_challenge_method"`
}

// AuthorizationCode is generated after successful user authentication and provider validation.
// needed for subsequent exchange for a token
type AuthorizationCode struct {
	Code        string    `json:"code"`
	SessionID   SessionID `json:"session_id"`
	UserID      UserID    `json:"user_id"`
	ClientID    ClientID  `json:"client_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       []Scope   `json:"scope"`
	Nonce       *Nonce    `json:"nonce"`
	PKCE        *PKCE     `json:"pkce"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Used        bool      `json:"used"`
}
