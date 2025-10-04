package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type JWKSRepository interface {
	GetCurrentJWKS(ctx context.Context) (*domain.JWKS, error)
	RotateKeys(ctx context.Context) error
}
