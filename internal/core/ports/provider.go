package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

// ProviderAuthenticator abstracts the interaction with external identity providers (e.g., Google, GitHub).
type ProviderAuthenticator interface {
	// ExchangeCodeForProfile takes an authorization code from the provider's callback
	// and returns a standardized ProviderProfile.
	ExchangeCodeForProfile(ctx context.Context, provider domain.Provider, code string) (*domain.ProviderProfile, error)
	// GetAuthURL returns the URL to redirect the user to for authentication.
	GetAuthURL(ctx context.Context, provider domain.Provider, state string) (string, error)
}
