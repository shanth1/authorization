// internal/ports/ports.go
package ports

import (
	"time"

	"github.com/shanth1/authorization/internal/domain"
)

// English comment: Ports (interfaces) for adapters, scalable to any impl.

// UserRepository port
type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByProvider(providerType string, sub string) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
}

// ClientRepository
type ClientRepository interface {
	FindByID(id string) (*domain.Client, error)
	Create(client *domain.Client) error
}

// SessionRepository
type SessionRepository interface {
	Create(session *domain.Session) error
	FindByID(id string) (*domain.Session, error)
	DeleteByUserID(userID string) error
}

// TokenRepository
type TokenRepository interface {
	Create(token *domain.Token) error
	FindByValue(value string) (*domain.Token, error)
	DeleteByValue(value string) error
	DeleteByUserID(userID string) error
}

// CachePort for JWKS, etc.
type CachePort interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
}

// ProviderPort - abstract for any provider
type ProviderPort interface {
	Authenticate(params map[string]string) (*domain.ProviderIdentity, error) // Returns identity or err
	GetAuthURL(clientID, state, nonce string) string                         // For external, but for internal empty
	HandleCallback(code string) (*domain.ProviderIdentity, error)            // For external
}
