package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/authorization/internal/core/domain"
)

type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (r *MemoryUserRepository) FindByProvider(ctx context.Context, providerType, providerExternalID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		for _, provider := range user.Providers {
			if provider.Type == providerType && provider.ExternalID == providerExternalID {
				return user, nil
			}
		}
	}

	return nil, errors.New("not found")
}

func (r *MemoryUserRepository) Upsert(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == "" {
		user.ID = generateUUID()
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	return nil
}

func generateUUID() string {
	return "user_" + uuid.New().String()
}
