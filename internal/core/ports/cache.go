// internal/core/ports/cache.go
package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type StateCache interface {
	SetState(ctx context.Context, key string, expiration time.Duration) error
	CheckState(ctx context.Context, key string) error
	DeleteState(ctx context.Context, key string) error
}

type SessionCache interface {
	SetSession(ctx context.Context, key string, tokens *domain.Tokens, expiration time.Duration) error
	GetSession(ctx context.Context, key string) (*domain.Tokens, error)
	DeleteSession(ctx context.Context, key string) error
}
