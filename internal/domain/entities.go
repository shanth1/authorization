// internal/domain/entities.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// English comment: Core domain entities, independent of tech.

// User entity
type User struct {
	ID        uuid.UUID
	Email     string
	Password  string // Hashed
	Providers []ProviderIdentity
}

// ProviderIdentity for merge
type ProviderIdentity struct {
	ProviderType string // "login_password", "google", etc.
	Sub          string // Subject ID from provider
	Claims       map[string]interface{}
}

// Client entity
type Client struct {
	ID           uuid.UUID
	Name         string
	Secret       string // Hashed if needed
	RedirectURIs []string
	Type         string // "public" or "confidential"
	Scopes       []string
	LogoutURI    string // For SSO logout
}

// Session for SSO
type Session struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ClientID      uuid.UUID
	State         string
	Nonce         string
	CodeChallenge string
	ExpiresAt     time.Time
}

// Token (code, access, refresh)
type Token struct {
	Value     string
	Type      string // "code", "access", "refresh"
	UserID    uuid.UUID
	ClientID  uuid.UUID
	Scope     string
	ExpiresAt time.Time
}
