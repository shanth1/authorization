package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type AuthorizationSessionRepository interface {
	Get(ctx context.Context, id domain.SessionID) (*domain.AuthorizationSession, error)
	Save(ctx context.Context, s *domain.AuthorizationSession, ttl time.Duration) error
	Delete(ctx context.Context, id domain.SessionID) error
}
