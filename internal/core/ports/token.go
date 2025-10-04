package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type TokenRepository interface {
	SaveAccess(ctx context.Context, t *domain.AccessToken, ttl time.Duration) error
	SaveRefresh(ctx context.Context, t *domain.RefreshToken, ttl time.Duration) error
	GetAccess(ctx context.Context, id domain.TokenID) (*domain.AccessToken, error)
	GetRefresh(ctx context.Context, id domain.TokenID) (*domain.RefreshToken, error)
	RevokeByClientUser(ctx context.Context, clientID domain.ClientID, userID *domain.UserID) error
	RevokeByJTI(ctx context.Context, jti domain.JTI) error
}
