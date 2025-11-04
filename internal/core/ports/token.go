package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type TokenRepository interface {
	// Access Token
	SaveAccess(ctx context.Context, t *domain.AccessToken, ttl time.Duration) error
	GetAccess(ctx context.Context, id domain.TokenID) (*domain.AccessToken, error)

	// Refresh Token
	SaveRefresh(ctx context.Context, t *domain.RefreshToken, ttl time.Duration) error
	GetRefresh(ctx context.Context, id domain.TokenID) (*domain.RefreshToken, error)
	GetRefreshByJTI(ctx context.Context, jti domain.JTI) (*domain.RefreshToken, error)

	// Revocation
	RevokeRefreshByJTI(ctx context.Context, jti domain.JTI) error
	RevokeAllByUserSession(ctx context.Context, sessionID domain.UserSessionID) error
	RevokeAllByClientUser(ctx context.Context, clientID domain.ClientID, userID domain.UserID) error
	RevokeAllByUser(ctx context.Context, userID domain.UserID) error

	// Listing
	ListRefreshByUser(ctx context.Context, userID domain.UserID) ([]*domain.RefreshToken, error)
	ListRefreshByUserSession(ctx context.Context, sessionID domain.UserSessionID) ([]*domain.RefreshToken, error)
}
