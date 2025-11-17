package service

import (
	"context"
	"fmt"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

// AuthenticationService manages user authentication.
type AuthenticationService struct {
	userRepo       ports.UserRepository
	sessionService *SessionService
	passwordHasher ports.PasswordHasher
	random         ports.Random
	clock          ports.Clock
}

func NewAuthenticationService(
	userRepo ports.UserRepository,
	sessionService *SessionService,
	passwordHasher ports.PasswordHasher,
	random ports.Random,
	clock ports.Clock,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo:       userRepo,
		sessionService: sessionService,
		passwordHasher: passwordHasher,
		random:         random,
		clock:          clock,
	}
}

// AuthenticateWithPassword authenticates the user by email/password
func (s *AuthenticationService) AuthenticateWithPassword(
	ctx context.Context,
	email string,
	password string,
	ipAddress string,
	userAgent string,
) (*domain.UserSession, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.Active {
		return nil, fmt.Errorf("user account is disabled")
	}

	if user.PasswordHash == nil {
		return nil, fmt.Errorf("password authentication not available")
	}

	valid := s.passwordHasher.Verify(password, *user.PasswordHash)
	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	session, err := s.sessionService.CreateSession(
		ctx,
		user.ID,
		"password",
		ipAddress,
		userAgent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// AuthenticateWithProvider authenticates through an external provider (Google, Telegram)
func (s *AuthenticationService) AuthenticateWithProvider(
	ctx context.Context,
	profile *domain.ProviderProfile,
	ipAddress string,
	userAgent string,
) (*domain.UserSession, error) {
	user, err := s.userRepo.GetByProviderID(ctx, profile.Provider, string(profile.Subject))

	// TODO: not found checking
	if err != nil {
		user, err = s.createUserFromProvider(ctx, profile)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	session, err := s.sessionService.CreateSession(
		ctx,
		user.ID,
		profile.Provider,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *AuthenticationService) createUserFromProvider(
	ctx context.Context,
	profile *domain.ProviderProfile,
) (*domain.User, error) {
	userID := domain.UserID(s.random.NewID())

	emails := []domain.Email{}
	if profile.Email != nil {
		emails = append(emails, *profile.Email)
	}

	user := &domain.User{
		ID:          userID,
		Active:      true,
		Emails:      emails,
		Username:    profile.Username,
		PictureURL:  profile.PictureURL,
		DisplayName: profile.DisplayName,
		CreatedAt:   s.clock.Now(),
	}

	err := s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	providerAccount := domain.ProviderAccount{
		UserID:     userID,
		Provider:   profile.Provider,
		Subject:    profile.Subject,
		Email:      profile.Email,
		Username:   profile.Username,
		RawProfile: profile.Raw,
		LinkedAt:   s.clock.Now(),
	}

	err = s.userRepo.LinkProvider(ctx, userID, providerAccount)
	if err != nil {
		return nil, err
	}

	return user, nil
}
