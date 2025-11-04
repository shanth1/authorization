package domain

import "time"

type Code string
type CodeChallengeMethod string // Method for generating code_challenge from code_verifier

const (
	CodeChallengeS256 CodeChallengeMethod = "S256"
)

// PKCE stores parameters for the Proof Key for Code Exchange mechanism (RFC 7636).
type PKCE struct {
	CodeChallenge       string              `json:"code_challenge"`
	CodeChallengeMethod CodeChallengeMethod `json:"code_challenge_method"`
}

// AuthorizationCode is generated after successful user authentication and provider validation.
// needed for subsequent exchange for a token
type AuthorizationCode struct {
	Code          Code          `json:"code"`
	RequestID     RequestID     `json:"request_id"` // ID of the original authorization request
	UserID        UserID        `json:"user_id"`
	UserSessionID UserSessionID `json:"user_session_id"`

	ClientID    ClientID  `json:"client_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       []Scope   `json:"scope"`
	Nonce       *Nonce    `json:"nonce"`
	PKCE        *PKCE     `json:"pkce"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"` // The code must have a very short lifetime (60 seconds)
	Used        bool      `json:"used"`
}
