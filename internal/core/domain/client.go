package domain

import "time"

type ClientID string
type Scope string
type GrantType string
type ResponseType string
type URI string

const (
	ResponseCode ResponseType = "code"

	ScopeOpenID  Scope = "openid"
	ScopeProfile Scope = "profile"
	ScopeEmail   Scope = "email"

	GrantAuthorizationCode GrantType = "authorization_code"
)

type Client struct {
	ID                   ClientID       `json:"id"`
	Name                 string         `json:"name"`
	Active               bool           `json:"active"`
	RedirectURIs         []URI          `json:"redirect_uris"`
	SecretHash           string         `json:"secret_hash"`
	BackchannelLogoutURI *URI           `json:"backchannel_logout_uri"`
	Scopes               []Scope        `json:"scopes"`         // [min] "openid"
	ResponseTypes        []ResponseType `json:"response_types"` // [default] "code"
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`

	ownerID                 UserID      // [optional] for user-created clients
	grantTypes              []GrantType // [optional] "authorization_code" | "refresh_token" | "implicit"
	TokenEndpointAuthMethod string      // [optional] "client_secret_basic" | "client_secret_post" | "none"
}
