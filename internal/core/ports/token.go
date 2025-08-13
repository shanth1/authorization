package ports

import "github.com/shanth1/authorization/internal/core/domain"

type TokenService interface {
	GenerateTokens(userID string) (*domain.Tokens, error)
}
