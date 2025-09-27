package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type AuthorizationCodeRepository interface {
	Save(ctx context.Context, code *domain.AuthorizationCode, ttl time.Duration) error
	Get(ctx context.Context, id string) (*domain.AuthorizationCode, error)
	MarkAsUsed(ctx context.Context, id string) error
}
