package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type AuthService interface {
	InitiateTelegramLogin(ctx context.Context) (string, error) // Returns login URL
	CompleteTelegramLogin(ctx context.Context, state string, telegramID int64, username string) error
	PollLoginStatus(ctx context.Context, state string) (*domain.Tokens, error)
}
