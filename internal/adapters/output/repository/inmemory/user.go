package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/authorization/internal/core/domain"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[int64]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[int64]*domain.User)}
}

func (r *InMemoryUserRepository) FindByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[telegramID]
	if !ok {
		return nil, fmt.Errorf("user not found") // Or a custom error type
	}
	return user, nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.TelegramID]; exists {
		return fmt.Errorf("user with telegram id %d already exists", user.TelegramID)
	}
	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	r.users[user.TelegramID] = user
	fmt.Printf("User created in-memory: %+v\n", user)
	return nil
}
