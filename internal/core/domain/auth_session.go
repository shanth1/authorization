package domain

import "time"

type RequestID string // request_id из URL
type State string     // CSRF protection
type Nonce string     // Replay attack protection
type Prompt string    // OIDC prompt parameter

const (
	PromptNone          Prompt = "none"           // Don't show UI (fail if not logged in)
	PromptLogin         Prompt = "login"          // Force the user to log in again
	PromptConsent       Prompt = "consent"        // Force the consent screen to appear
	PromptSelectAccount Prompt = "select_account" // Show account selection
)

// AuthorizationRequest stores the parameters of the original authorization request
// is created in the first step and is used for verification in subsequent steps
// ensures security before code is released
type AuthorizationRequest struct {
	ID          RequestID `json:"id"`
	ClientID    ClientID  `json:"client_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       []Scope   `json:"scope"`
	State       State     `json:"state"`
	Nonce       *Nonce    `json:"nonce"` // Associates a client session with an ID Token to protect against replay attacks
	PKCE        *PKCE     `json:"pkce"`
	Prompt      *Prompt   `json:"prompt"` // UI behavior (none, login, agreement, select_account)
	Provider    Provider  `json:"provider"`
	MaxAge      *int      `json:"max_age"` // ?? Maximum authentication age in seconds

	// Request Processing Status
	UserSessionID  *UserSessionID `json:"user_session_id"` // Filled in after authentication
	ConsentGranted bool           `json:"consent_granted"` // Filled out after consent (if required)

	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"` // The session should have a short lifetime (3-10 minutes)

}
