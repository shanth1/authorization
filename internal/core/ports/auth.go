package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type EmailSender interface {
	SendVerificationCode(ctx context.Context, to string, code string, ttl time.Duration) error
}

type VerificationCodeStore interface {
	Save(ctx context.Context, email string, code string, ttl time.Duration) error
	VerifyAndConsume(ctx context.Context, email string, code string) (bool, error)
}

type LocalAuth interface {
	Register(ctx context.Context, email string, password string) (*domain.User, error)
	LoginWithPassword(ctx context.Context, email string, password string) (*domain.User, error)
	BeginEmailVerification(ctx context.Context, email string) error
	CompleteEmailVerification(ctx context.Context, email string, code string) (*domain.User, error)
}
