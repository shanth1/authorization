package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type DiscoveryRepository interface {
	Get(ctx context.Context) (*domain.Discovery, error)
	Save(ctx context.Context, d *domain.Discovery) error
}
