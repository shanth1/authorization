package domain

import "time"

type UserID string

type User struct {
	ID               UserID            `json:"id"`
	Active           bool              `json:"active"`
	Emails           []Email           `json:"emails"`
	Username         *string           `json:"username"`
	PictureURL       *string           `json:"picture_url"`
	DisplayName      *string           `json:"display_name"`
	PasswordHash     *string           `json:"password_hash"`
	Profile          map[string]any    `json:"profile"` // description, locale, etc.
	ProviderAccounts []ProviderAccount `json:"provider_accounts"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        *time.Time        `json:"updated_at"`
}
