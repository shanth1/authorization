package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

// UserRepository управляет пользователями
type UserRepository interface {
	// FindByID ищет пользователя по внутреннему ID
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByProvider ищет пользователя по данным провайдера
	FindByProvider(
		ctx context.Context,
		providerType string,
		externalID string,
	) (*domain.User, error)

	// Upsert создает или обновляет пользователя
	Upsert(ctx context.Context, user *domain.User) error
}

// OIDCClientRepository управляет OIDC-клиентами
type OIDCClientRepository interface {
	// Get возвращает клиента по ID
	Get(ctx context.Context, clientID string) (*domain.OIDCClient, error)
}

// AuthCodeRepository управляет кодами авторизации
type AuthCodeRepository interface {
	// Store сохраняет код авторизации
	Store(
		ctx context.Context,
		code string,
		data *domain.AuthCodeData,
	) error

	// Retrieve возвращает данные по коду и удаляет его
	Retrieve(
		ctx context.Context,
		code string,
	) (*domain.AuthCodeData, error)
}
