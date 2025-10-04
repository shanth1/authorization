package v1

import "github.com/gin-gonic/gin"

type Handlers struct {
	Discovery *DiscoveryHandler
	JWKS      *JWKSHandler
	OIDC      *OIDCHandler
	User      *UserHandler
	Provider  *ProviderHandler
	Admin     *AdminHandler
}

type Middlewares struct {
	UserSession gin.HandlerFunc
	BearerToken gin.HandlerFunc
	Admin       gin.HandlerFunc
}

type (
	DiscoveryHandler struct{ GetConfiguration gin.HandlerFunc }
	JWKSHandler      struct{ GetKeys gin.HandlerFunc }
	OIDCHandler      struct {
		Authorize        gin.HandlerFunc
		Token            gin.HandlerFunc
		Revoke           gin.HandlerFunc
		Introspect       gin.HandlerFunc
		EndSession       gin.HandlerFunc
		UserInfo         gin.HandlerFunc
		ClientWideLogout gin.HandlerFunc
	}
	UserHandler struct {
		Register         gin.HandlerFunc
		LoginPassword    gin.HandlerFunc
		GetMyProfile     gin.HandlerFunc
		GetMySessions    gin.HandlerFunc
		TerminateSession gin.HandlerFunc
		GlobalLogout     gin.HandlerFunc
	}
	ProviderHandler struct {
		List             gin.HandlerFunc
		GoogleCallback   gin.HandlerFunc
		TelegramCallback gin.HandlerFunc
	}
	AdminHandler struct {
		ListClients          gin.HandlerFunc
		CreateClient         gin.HandlerFunc
		GetClient            gin.HandlerFunc
		UpdateClient         gin.HandlerFunc
		DeleteClient         gin.HandlerFunc
		ListUsers            gin.HandlerFunc
		CreateUser           gin.HandlerFunc
		GetUser              gin.HandlerFunc
		UpdateUser           gin.HandlerFunc
		DeleteUser           gin.HandlerFunc
		TerminateUserSession gin.HandlerFunc
		GlobalLogoutUser     gin.HandlerFunc
		GetKeys              gin.HandlerFunc
		RotateKeys           gin.HandlerFunc
		GetAuditLogs         gin.HandlerFunc
	}
)
