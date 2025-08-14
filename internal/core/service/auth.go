package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

type AuthService struct {
	userRepo  ports.UserRepository
	oidcRepo  ports.OIDCClientRepository
	providers map[string]ports.AuthProvider
}

func NewAuthService(
	userRepo ports.UserRepository,
	oidcRepo ports.OIDCClientRepository,
	providers []ports.AuthProvider,
) *AuthService {
	providerMap := make(map[string]ports.AuthProvider)
	for _, p := range providers {
		providerMap[p.Type()] = p
	}

	return &AuthService{
		userRepo:  userRepo,
		oidcRepo:  oidcRepo,
		providers: providerMap,
	}
}

// HandleOAuthCallback обрабатывает OAuth callback
func (s *AuthService) HandleOAuthCallback(
	ctx context.Context,
	providerType string,
	code string,
) (*domain.User, error) {
	provider, exists := s.providers[providerType]
	if !exists {
		return nil, errors.New("provider not supported")
	}

	return provider.Exchange(ctx, code)
}

// CompleteAuth завершает аутентификацию
func (s *AuthService) CompleteAuth(
	ctx context.Context,
	providerType string,
	externalID string,
) (*domain.User, error) {
	// Находим пользователя по идентификатору провайдера
	user, err := s.userRepo.FindByProvider(ctx, providerType, externalID)
	if err != nil {
		// Создаем нового пользователя, если не найден
		user = &domain.User{
			ID: uuid.New().String(),
			Providers: []domain.Provider{{
				Type:       providerType,
				ExternalID: externalID,
			}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// Обновляем информацию о пользователе
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Upsert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
