package domain

import "time"

// UserSessionID is the unique identifier for a user's SSO session.
// This is the value stored in the user's cookie.
type UserSessionID string

// UserSession represents a user's SSO session in the Authorization Server
// This is a long-lived session (up to 30 days), which allows the user
// to avoid having to re-enter credentials when logging in to different clients
type UserSession struct {
	ID         UserSessionID `json:"id"`
	UserID     UserID        `json:"user_id"`
	Provider   Provider      `json:"provider"`
	IPAddress  string        `json:"ip_address"`
	UserAgent  string        `json:"user_agent"` // Browser User-Agent
	CreatedAt  time.Time     `json:"created_at"`
	ExpiresAt  time.Time     `json:"expires_at"`
	LastUsedAt time.Time     `json:"last_used_at"`
}
