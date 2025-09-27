package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type JTIRepository interface {
	MarkUsed(ctx context.Context, jti domain.JTI, ttl time.Duration) error
	IsUsed(ctx context.Context, jti domain.JTI) (bool, error)
}
