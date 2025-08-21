package adapters

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/authorization/internal/domain"
	"github.com/shanth1/authorization/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

// English comment: In-memory adapters for repositories and cache. These are thread-safe using sync.Map and suitable for development/testing. For production, replace with DB adapters.

// InMemoryUserRepository
type InMemoryUserRepository struct {
	users sync.Map // key: user.ID.String(), value: *domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{}
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
	var found *domain.User
	r.users.Range(func(key, value interface{}) bool {
		user := value.(*domain.User)
		if user.Email == email {
			found = user
			return false
		}
		return true
	})
	if found == nil {
		return nil, errors.New("user not found")
	}
	return found, nil
}

func (r *InMemoryUserRepository) FindByProvider(providerType string, sub string) (*domain.User, error) {
	var found *domain.User
	r.users.Range(func(key, value interface{}) bool {
		user := value.(*domain.User)
		for _, p := range user.Providers {
			if p.ProviderType == providerType && p.Sub == sub {
				found = user
				return false
			}
		}
		return true
	})
	if found == nil {
		return nil, errors.New("user not found by provider")
	}
	return found, nil
}

func (r *InMemoryUserRepository) Create(user *domain.User) error {
	r.users.Store(user.ID.String(), user)
	return nil
}

func (r *InMemoryUserRepository) Update(user *domain.User) error {
	_, ok := r.users.Load(user.ID.String())
	if !ok {
		return errors.New("user not found for update")
	}
	r.users.Store(user.ID.String(), user)
	return nil
}

// InMemoryClientRepository
type InMemoryClientRepository struct {
	clients sync.Map // key: client.ID.String(), value: *domain.Client
}

func NewInMemoryClientRepository() *InMemoryClientRepository {
	return &InMemoryClientRepository{}
}

func (r *InMemoryClientRepository) FindByID(id string) (*domain.Client, error) {
	v, ok := r.clients.Load(id)
	if !ok {
		return nil, errors.New("client not found")
	}
	return v.(*domain.Client), nil
}

func (r *InMemoryClientRepository) Create(client *domain.Client) error {
	r.clients.Store(client.ID.String(), client)
	return nil
}

// InMemorySessionRepository
type InMemorySessionRepository struct {
	sessions sync.Map // key: session.ID.String(), value: *domain.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{}
}

func (r *InMemorySessionRepository) Create(session *domain.Session) error {
	r.sessions.Store(session.ID.String(), session)
	return nil
}

func (r *InMemorySessionRepository) FindByID(id string) (*domain.Session, error) {
	v, ok := r.sessions.Load(id)
	if !ok {
		return nil, errors.New("session not found")
	}
	session := v.(*domain.Session)
	if time.Now().After(session.ExpiresAt) {
		r.sessions.Delete(id) // Clean up expired
		return nil, errors.New("session expired")
	}
	return session, nil
}

func (r *InMemorySessionRepository) DeleteByUserID(userID string) error {
	r.sessions.Range(func(key, value interface{}) bool {
		session := value.(*domain.Session)
		if session.UserID.String() == userID {
			r.sessions.Delete(key)
		}
		return true
	})
	return nil
}

// InMemoryTokenRepository
type InMemoryTokenRepository struct {
	tokens sync.Map // key: token.Value, value: *domain.Token
}

func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{}
}

func (r *InMemoryTokenRepository) Create(token *domain.Token) error {
	r.tokens.Store(token.Value, token)
	return nil
}

func (r *InMemoryTokenRepository) FindByValue(value string) (*domain.Token, error) {
	v, ok := r.tokens.Load(value)
	if !ok {
		return nil, errors.New("token not found")
	}
	token := v.(*domain.Token)
	if time.Now().After(token.ExpiresAt) {
		r.tokens.Delete(value) // Clean up expired
		return nil, errors.New("token expired")
	}
	return token, nil
}

func (r *InMemoryTokenRepository) DeleteByValue(value string) error {
	r.tokens.Delete(value)
	return nil
}

func (r *InMemoryTokenRepository) DeleteByUserID(userID string) error {
	r.tokens.Range(func(key, value interface{}) bool {
		token := value.(*domain.Token)
		if token.UserID.String() == userID {
			r.tokens.Delete(key)
		}
		return true
	})
	return nil
}

// InMemoryCache (already provided, but repeated for completeness)
type InMemoryCache struct {
	cache sync.Map
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.cache.Store(key, value)
	// For TTL, spawn a goroutine to delete after ttl (simple in-memory impl)
	go func() {
		<-time.After(ttl)
		c.cache.Delete(key)
	}()
	return nil
}

func (c *InMemoryCache) Get(key string) (interface{}, error) {
	v, ok := c.cache.Load(key)
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

// LoginPasswordProvider (already provided, no changes)
type LoginPasswordProvider struct {
	userRepo ports.UserRepository
}

func NewLoginPasswordProvider(userRepo ports.UserRepository) *LoginPasswordProvider {
	return &LoginPasswordProvider{userRepo: userRepo}
}

func (p *LoginPasswordProvider) Authenticate(params map[string]string) (*domain.ProviderIdentity, error) {
	login, ok := params["login"]
	if !ok {
		return nil, errors.New("missing login")
	}
	password, ok := params["password"]
	if !ok {
		return nil, errors.New("missing password")
	}

	user, err := p.userRepo.FindByEmail(login)
	if err != nil {
		// Create new user if not found (for demo; in prod, separate registration)
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user = &domain.User{
			ID:       uuid.New(),
			Email:    login,
			Password: string(hashed),
		}
		if err := p.userRepo.Create(user); err != nil {
			return nil, err
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, errors.New("invalid password")
		}
	}

	return &domain.ProviderIdentity{
		ProviderType: "login_password",
		Sub:          user.ID.String(),
		Claims:       map[string]interface{}{"email": user.Email},
	}, nil
}

func (p *LoginPasswordProvider) GetAuthURL(_, _, _ string) string {
	return "" // Internal, no redirect
}

func (p *LoginPasswordProvider) HandleCallback(_ string) (*domain.ProviderIdentity, error) {
	return nil, nil // Not used for internal
}
