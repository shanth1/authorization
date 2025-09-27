package domain

import "time"

type ExternalID string
type Provider string // provider type (google, github, telegram, etc.)

// ProviderProfile contains temporary data received from the provider at the login stage
type ProviderProfile struct {
	Provider    Provider       `json:"provider"`
	Subject     ExternalID     `json:"subject"` // (e.g., Google sub, GitHub id, Telegram id)
	Email       *Email         `json:"email"`
	Username    *string        `json:"username"`
	PictureURL  *string        `json:"picture_url"`
	DisplayName *string        `json:"display_name"`
	Raw         map[string]any `jspn:"raw"` // raw data received from the provider
}

// ProviderAccount contains persistent data that links the user to a specific provider
type ProviderAccount struct {
	UserID     UserID         `json:"user_id"`
	Provider   Provider       `json:"provider"`
	Subject    ExternalID     `json:"subject"`
	Email      *Email         `json:"email"`
	Username   *string        `json:"username"`
	RawProfile map[string]any `json:"raw_profile"` // JSON of the last received fields (name, picture, locale, etc.)
	LinkedAt   time.Time      `json:"linked_at"`

	issuer string // [optional] provider url
}
