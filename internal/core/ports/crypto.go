package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type Clock interface {
	Now() time.Time
}

type Random interface {
	NewID() string
	Bytes(n int) ([]byte, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

type PKCEService interface {
	Verify(codeVerifier string, challenge string, method domain.CodeChallengeMethod) bool
}

type StateNonceService interface {
	GenerateState() domain.State
	GenerateNonce() domain.Nonce
}

type JWTSigner interface {
	SignIDToken(ctx context.Context, claims domain.IDTokenClaims) (string, error)
	SignLogoutToken(ctx context.Context, claims domain.LogoutTokenClaims) (string, error)
}

type JWTVerifier interface {
	VerifyIDToken(ctx context.Context, token string, expectedNonce *domain.Nonce, audience string, issuer string) (*domain.IDTokenClaims, error)
}

type JWKSManager interface {
	CurrentJWKS(ctx context.Context) (*domain.JWKS, error)
}
