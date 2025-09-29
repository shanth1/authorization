package inmemory

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/gotools/errs"
)

type JWKSRepository interface {
	GetCurrentJWKS(ctx context.Context) (*domain.JWKS, error)
	RotateKeys(ctx context.Context) error
}
type jwksRepo struct {
	jwks *domain.JWKS
}

func NewJWKSRepo() *jwksRepo {
	return &jwksRepo{
		jwks: &domain.JWKS{
			Keys: []domain.JWK{},
		},
	}
}

func (r *jwksRepo) GetCurrentJWKS(ctx context.Context) (*domain.JWKS, error) {
	return r.jwks, nil
}

func (r *jwksRepo) RotateKeys(ctx context.Context) error {
	return errs.ErrNotImplemented
}
