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
	GrantRefreshToken      GrantType = "refresh_token"
)

type Client struct {
	ID                   ClientID       `json:"id"`
	Name                 string         `json:"name"`
	Active               bool           `json:"active"`
	RedirectURIs         []URI          `json:"redirect_uris"` // List of allowed URIs for redirection after authorization
	SecretHash           string         `json:"secret_hash"`
	BackchannelLogoutURI *URI           `json:"backchannel_logout_uri"` // URI for Back-Channel Logout. A POST-request with a logout token is sent to this address
	Scopes               []Scope        `json:"scopes"`                 // [min] "openid"
	ResponseTypes        []ResponseType `json:"response_types"`         // authorization flow type. [default] "code"
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            *time.Time     `json:"updated_at"`
	GrantTypes           []GrantType    `json:"grant_types"`

	ownerID                 UserID // [optional] for user-created clients
	tokenEndpointAuthMethod string // Client authentication method on the /token endpoint. [optional] "client_secret_basic" | "client_secret_post" | "none"
}
