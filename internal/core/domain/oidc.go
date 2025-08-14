package domain

import "time"

// OIDCClient представляет зарегистрированное клиентское приложение в IdP
// TODO: rename to app?
type OIDCClient struct {
	ID           string   `json:"id"`            // Идентификатор клиента
	Secret       string   `json:"secret"`        // Секрет клиента
	RedirectURIs []string `json:"redirect_uris"` // Доверенные URI для редиректов
	Name         string   `json:"name"`          // Человекочитаемое название
}

// AuthorizationRequest - запрос на авторизацию (RFC 6749 §4.1.1)
type AuthorizationRequest struct {
	ResponseType string `form:"response_type"` // Тип ответа (например "code")
	ClientID     string `form:"client_id"`     // Идентификатор клиента
	RedirectURI  string `form:"redirect_uri"`  // URI для редиректа после авторизации
	Scope        string `form:"scope"`         // Запрашиваемые разрешения
	State        string `form:"state"`         // CSRF-токен клиента
	Nonce        string `form:"nonce"`         // Уникальный токен для привязки сессии
}

// AuthCodeData хранит данные сессии авторизации
type AuthCodeData struct {
	ClientID    string    // Идентификатор клиента
	RedirectURI string    // URI для редиректа
	UserID      string    // Идентификатор пользователя
	Nonce       string    // Уникальный токен сессии
	ExpiresAt   time.Time // Время истечения
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	Nickname      string `json:"nickname,omitempty"`
	PreferredName string `json:"preferred_username,omitempty"`
	Profile       string `json:"profile,omitempty"`
	Picture       string `json:"picture,omitempty"`
	Website       string `json:"website,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Gender        string `json:"gender,omitempty"`
	Birthdate     string `json:"birthdate,omitempty"`
	Zoneinfo      string `json:"zoneinfo,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Phone         string `json:"phone_number,omitempty"`
	PhoneVerified bool   `json:"phone_number_verified,omitempty"`
	Address       string `json:"address,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
}
