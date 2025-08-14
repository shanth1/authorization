package domain

import "time"

// User - основной объект пользователя IdP
type User struct {
	ID           string     `json:"id"`            // Внутренний UUID
	PrimaryEmail string     `json:"primary_email"` // Основной email (верифицированный)
	Providers    []Provider `json:"providers"`     // Подключенные провайдеры
	CreatedAt    time.Time  `json:"created_at"`    // Время создания
	UpdatedAt    time.Time  `json:"updated_at"`    // Время последнего обновления
	PasswordHash *string    `json:"-"`             // not used
}

// Provider представляет подключенный OAuth-провайдер
type Provider struct {
	Type         string    `json:"type"`          // Тип провайдера (telegram, google и т.д.)
	ExternalID   string    `json:"external_id"`   // ID пользователя в системе провайдера
	Email        string    `json:"email"`         // Email из провайдера
	DisplayName  string    `json:"display_name"`  // Отображаемое имя
	AvatarURL    string    `json:"avatar_url"`    // Ссылка на аватар
	AccessToken  string    `json:"access_token"`  // Токен доступа (опционально)
	RefreshToken string    `json:"refresh_token"` // Refresh-токен (опционально)
	ExpiresAt    time.Time `json:"expires_at"`    // Время истечения токена
	ConnectedAt  time.Time `json:"connected_at"`  // Время подключения
}
