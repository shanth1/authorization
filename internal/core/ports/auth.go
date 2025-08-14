package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type AuthService interface {
	HandleOAuthCallback(
		ctx context.Context,
		providerType string,
		code string,
	) (*domain.User, error)

	CompleteAuth(
		ctx context.Context,
		providerType string,
		externalID string,
	) (*domain.User, error)
}

// AuthProvider интерфейс для внешних провайдеров (Telegram, Google и т.д.)
type AuthProvider interface {
	BuildAuthURL(state string) string                                // BuildAuthURL создает URL для аутентификации через провайдера
	Exchange(ctx context.Context, code string) (*domain.User, error) // Exchange обменивает код на данные пользователя
	Type() string
}

// OIDCService интерфейс для ядра OIDC
type OIDCService interface {
	// HandleAuthorize обрабатывает OIDC-запрос авторизации
	HandleAuthorize(
		ctx context.Context,
		req *domain.AuthorizationRequest,
	) (redirectURL string, err error)

	// ExchangeCode обменивает код авторизации на токены
	ExchangeCode(
		ctx context.Context,
		code string,
	) (*domain.Tokens, error)

	// GetUserInfo возвращает информацию о пользователе
	GetUserInfo(
		ctx context.Context,
		accessToken string,
	) (*domain.User, error)
}
