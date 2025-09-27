package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type UserRepository interface {
	Get(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetByProviderID(ctx context.Context, provider domain.Provider, sub string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, u *domain.User) error
	LinkProvider(ctx context.Context, userID domain.UserID, acc domain.ProviderAccount) error
}
