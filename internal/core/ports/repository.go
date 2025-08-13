package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type UserRepository interface {
	FindByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}
