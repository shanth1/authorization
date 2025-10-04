package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/gotools/errs"
)

type authCodeRepo struct {
	mu    sync.RWMutex
	codes map[domain.Code]*domain.AuthorizationCode
}

func NewAuthCodeRepo() *authCodeRepo {
	return &authCodeRepo{
		codes: make(map[domain.Code]*domain.AuthorizationCode),
	}
}

func (r *authCodeRepo) Get(ctx context.Context, id domain.Code) (*domain.AuthorizationCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	code, exists := r.codes[id]
	if !exists {
		return nil, errs.ErrNotFound
	}

	codeCopy := *code
	return &codeCopy, nil
}

func (r *authCodeRepo) Save(ctx context.Context, c *domain.AuthorizationCode, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if c.Code == "" {
		return errors.New("emtpy id")
	}

	codeToSave := *c
	r.codes[c.Code] = &codeToSave

	return nil
}

func (r *authCodeRepo) MarkAsUsed(ctx context.Context, id domain.Code) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	authCode, exists := r.codes[id]
	if !exists {
		return errs.ErrNotFound
	}

	authCode.Used = true

	return nil
}
