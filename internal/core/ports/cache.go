package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

// StateCache управляет OAuth state-параметрами
type StateCache interface {
	// Store сохраняет state с данными авторизации
	Store(
		ctx context.Context,
		state string,
		data *domain.AuthorizationRequest,
		exp time.Duration,
	) error

	// Verify проверяет state и возвращает данные
	Verify(
		ctx context.Context,
		state string,
	) (*domain.AuthorizationRequest, error)
}

// SessionCache управляет пользовательскими сессиями
type SessionCache interface {
	// SetTokens сохраняет токены сессии
	SetTokens(
		ctx context.Context,
		sessionID string,
		tokens *domain.Tokens,
		exp time.Duration,
	) error

	// GetTokens возвращает токены сессии
	GetTokens(
		ctx context.Context,
		sessionID string,
	) (*domain.Tokens, error)
}
