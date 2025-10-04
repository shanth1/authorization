package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/gotools/errs"
)

type userRepo struct {
	mu            sync.RWMutex
	users         map[domain.UserID]*domain.User
	emailIndex    map[string]domain.UserID
	providerIndex map[string]domain.UserID
}

func NewUserRepo() *userRepo {
	return &userRepo{
		users:         make(map[domain.UserID]*domain.User),
		emailIndex:    make(map[string]domain.UserID),
		providerIndex: make(map[string]domain.UserID),
	}
}

func (r *userRepo) Get(ctx context.Context, id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errs.ErrNotFound
	}

	userCopy := *user
	return &userCopy, nil
}

func (r *userRepo) GetByProviderID(ctx context.Context, provider domain.Provider, sub string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	indexKey := fmt.Sprintf("%s:%s", provider, sub)
	userID, exists := r.providerIndex[indexKey]
	if !exists {
		return nil, errs.ErrNotFound
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, errs.ErrNotFound
	}

	userCopy := *user
	return &userCopy, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.emailIndex[email]
	if !exists {
		return nil, errs.ErrNotFound
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, errs.ErrNotFound
	}

	userCopy := *user
	return &userCopy, nil
}

func (r *userRepo) Save(ctx context.Context, u *domain.User) error {
	if u.ID == "" {
		return errors.New("empty id")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if existingUser, exists := r.users[u.ID]; exists {
		r.removeIndexes(existingUser)
	}

	userToSave := *u
	r.users[u.ID] = &userToSave

	r.addIndexes(&userToSave)

	return nil
}

func (r *userRepo) LinkProvider(ctx context.Context, userID domain.UserID, acc domain.ProviderAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return errs.ErrNotFound
	}

	indexKey := fmt.Sprintf("%s:%s", acc.Provider, acc.Subject)
	if _, exists := r.providerIndex[indexKey]; exists {
		return nil
	}

	user.ProviderAccounts = append(user.ProviderAccounts, acc)
	now := time.Now()
	user.UpdatedAt = &now

	r.providerIndex[indexKey] = userID

	return nil
}

func (r *userRepo) removeIndexes(u *domain.User) {
	for _, email := range u.Emails {
		delete(r.emailIndex, email.Address)
	}
	for _, acc := range u.ProviderAccounts {
		indexKey := fmt.Sprintf("%s:%s", acc.Provider, acc.Subject)
		delete(r.providerIndex, indexKey)
	}
}

func (r *userRepo) addIndexes(u *domain.User) {
	for _, email := range u.Emails {
		r.emailIndex[email.Address] = u.ID
	}
	for _, acc := range u.ProviderAccounts {
		indexKey := fmt.Sprintf("%s:%s", acc.Provider, acc.Subject)
		r.providerIndex[indexKey] = u.ID
	}
}
