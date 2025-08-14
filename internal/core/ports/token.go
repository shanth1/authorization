package ports

import "github.com/shanth1/authorization/internal/core/domain"

// TokenService управляет JWT-токенами
type TokenService interface {
	// GenerateTokens создает набор токенов
	GenerateTokens(
		user *domain.User,
		client *domain.OIDCClient,
		nonce string,
	) (*domain.Tokens, error)

	// ParseAccessToken валидирует и парсит access token
	ParseAccessToken(token string) (*domain.TokenClaims, error)

	// ParseIDToken валидирует и парсит ID token
	ParseIDToken(token string) (*domain.TokenClaims, error)
}
