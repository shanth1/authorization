package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/gotools/errs"
)

type authSessionRepo struct {
	mu       sync.RWMutex
	sessions map[domain.SessionID]*domain.AuthorizationSession
	stopChan chan struct{}
}

func NewAuthSessionRepo(cleanupInterval time.Duration) *authSessionRepo {
	repo := &authSessionRepo{
		sessions: make(map[domain.SessionID]*domain.AuthorizationSession),
		stopChan: make(chan struct{}),
	}

	return repo
}

func (r *authSessionRepo) Get(ctx context.Context, id domain.SessionID) (*domain.AuthorizationSession, error) {
	r.mu.RLock()
	session, exists := r.sessions[id]
	r.mu.RUnlock()

	if !exists {
		return nil, errs.ErrNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		r.mu.Lock()
		delete(r.sessions, id)
		r.mu.Unlock()
		return nil, errs.ErrNotFound
	}

	sessionCopy := *session
	return &sessionCopy, nil
}

func (r *authSessionRepo) Save(ctx context.Context, s *domain.AuthorizationSession, ttl time.Duration) error {
	if s.ID == "" {
		return errors.New("empty id")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	s.CreatedAt = time.Now()
	s.ExpiresAt = time.Now().Add(ttl)

	r.sessions[s.ID] = s

	return nil
}

func (r *authSessionRepo) Delete(ctx context.Context, id domain.SessionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, id)

	return nil
}
