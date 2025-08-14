package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shanth1/authorization/internal/core/domain"
)

type JWTService struct {
	signingKey []byte
}

func NewJWTService(key string) *JWTService {
	return &JWTService{signingKey: []byte(key)}
}

func (s *JWTService) GenerateTokens(
	user *domain.User,
	client *domain.OIDCClient,
	nonce string,
) (*domain.Tokens, error) {
	now := time.Now()
	expiresIn := int64(3600) // 1 час

	// Access Token
	accessClaims := &domain.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:   user.ID,
		ClientID: client.ID,
		Scopes:   []string{"openid", "profile"},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString(s.signingKey)
	if err != nil {
		return nil, err
	}

	// ID Token
	idClaims := &domain.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:   user.ID,
		ClientID: client.ID,
		Nonce:    nonce,
	}

	idToken := jwt.NewWithClaims(jwt.SigningMethodHS256, idClaims)
	idSigned, err := idToken.SignedString(s.signingKey)
	if err != nil {
		return nil, err
	}

	return &domain.Tokens{
		AccessToken: accessSigned,
		IDToken:     idSigned,
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *JWTService) ParseAccessToken(token string) (*domain.TokenClaims, error) {
	// Реализация валидации токена
	return nil, errors.New("not implemeted")
}
