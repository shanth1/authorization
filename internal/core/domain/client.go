package domain

import "time"

type ClientID string
type Scope string
type GrantType string
type ResponseType string
type URI string

type ClientType string
type TokenEndpointAuthMethod string

const (
	ResponseCode ResponseType = "code"

	ScopeOpenID  Scope = "openid"
	ScopeProfile Scope = "profile"
	ScopeEmail   Scope = "email"

	GrantAuthorizationCode GrantType = "authorization_code"
	GrantRefreshToken      GrantType = "refresh_token"

	ClientTypeFirstParty ClientType = "first_party" // Own clients (portal, SSO)
	ClientTypeThirdParty ClientType = "third_party" // External clients

	AuthMethodClientSecretBasic TokenEndpointAuthMethod = "client_secret_basic" // Basic auth
	AuthMethodClientSecretPost  TokenEndpointAuthMethod = "client_secret_post"  // In the body of the request
	AuthMethodNone              TokenEndpointAuthMethod = "none"                // Public clients (PKCE required)
)

type Client struct {
	ID                   ClientID       `json:"id"`
	Name                 string         `json:"name"`
	Type                 ClientType     `json:"type"`
	Active               bool           `json:"active"`
	RedirectURIs         []URI          `json:"redirect_uris"` // List of allowed URIs for redirection after authorization
	SecretHash           string         `json:"secret_hash"`
	BackchannelLogoutURI *URI           `json:"backchannel_logout_uri"` // URI for Back-Channel Logout. A POST-request with a logout token is sent to this address
	Scopes               []Scope        `json:"scopes"`                 // [min] "openid"
	ResponseTypes        []ResponseType `json:"response_types"`         // authorization flow type. [default] "code"
	GrantTypes           []GrantType    `json:"grant_types"`

	AllowedProviders []Provider `json:"allowed_providers"` // ?? ["password", "google", "telegram"]
	RequiresConsent  bool       `json:"requires_consent"`  // Whether a consent screen is required (usually true for third_party)

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	OwnerID                 UserID // ?? [optional] for user-created clients
	TokenEndpointAuthMethod string // ?? Client authentication method on the /token endpoint. [optional] "client_secret_basic" | "client_secret_post" | "none"
}
