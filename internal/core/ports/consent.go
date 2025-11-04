package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

// ConsentRepository manages user consents
type ConsentRepository interface {
	// Save saves or updates consent
	Save(ctx context.Context, consent *domain.ConsentGrant) error

	// Get obtains consent for the user-client pair
	Get(ctx context.Context, userID domain.UserID, clientID domain.ClientID) (*domain.ConsentGrant, error)

	// Exists checks for valid consent
	Exists(ctx context.Context, userID domain.UserID, clientID domain.ClientID, scopes []domain.Scope) (bool, error)

	// Revoke revokes consent
	Revoke(ctx context.Context, userID domain.UserID, clientID domain.ClientID) error

	// ListByUser obtains all user consents
	ListByUser(ctx context.Context, userID domain.UserID) ([]*domain.ConsentGrant, error)
}
