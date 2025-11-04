package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

// AuthorizationRequestRepository manages authorization requests (request_id)
type AuthorizationRequestRepository interface {
	Save(ctx context.Context, req *domain.AuthorizationRequest, ttl time.Duration) error
	Get(ctx context.Context, id domain.RequestID) (*domain.AuthorizationRequest, error)

	// Update refreshes the request (after authentication/consent)
	Update(ctx context.Context, req *domain.AuthorizationRequest) error

	// Delete deletes the request after issuing the code
	Delete(ctx context.Context, id domain.RequestID) error
}

// UserSessionRepository manages the user's SSO session.
type UserSessionRepository interface {
	Create(ctx context.Context, session *domain.UserSession, ttl time.Duration) error
	Get(ctx context.Context, id domain.UserSessionID) (*domain.UserSession, error)
	UpdateLastUsed(ctx context.Context, id domain.UserSessionID) error
	Delete(ctx context.Context, id domain.UserSessionID) error
	DeleteAllByUser(ctx context.Context, userID domain.UserID) error
	ListByUser(ctx context.Context, userID domain.UserID) ([]*domain.UserSession, error)
}
