package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type ClientRepository interface {
	Get(ctx context.Context, id domain.ClientID) (*domain.Client, error)
	Save(ctx context.Context, c *domain.Client) error
}
