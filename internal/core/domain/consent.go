package domain

import "time"

// ConsentGrant represents a user's consent to access their data
// for a specific client with specific scopes
type ConsentGrant struct {
	ID        string     `json:"id"`
	UserID    UserID     `json:"user_id"`
	ClientID  ClientID   `json:"client_id"`
	Scopes    []Scope    `json:"scopes"`
	GrantedAt time.Time  `json:"granted_at"`
	ExpiresAt *time.Time `json:"expires_at"` // nil = indefinitely
}

// HasScope checks whether the scope is included in the consent
func (c *ConsentGrant) HasScope(scope Scope) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// HasAllScopes checks whether all requested scopes are enabled.
func (c *ConsentGrant) HasAllScopes(requestedScopes []Scope) bool {
	for _, requested := range requestedScopes {
		if !c.HasScope(requested) {
			return false
		}
	}
	return true
}
