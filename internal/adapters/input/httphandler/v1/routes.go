package v1

import "github.com/gin-gonic/gin"

type v1Router struct {
	path        string
	handlers    *Handlers
	middlewares *Middlewares
}

func NewRouter() *v1Router {
	return &v1Router{
		path:        "/api/v1",
		handlers:    &Handlers{},
		middlewares: &Middlewares{},
	}
}

// TODO:
func (router *v1Router) SetupV1Routes(rg *gin.RouterGroup) {
	// =========================================================================
	// --- Public OIDC & OAuth2 Endpoints ---
	// These endpoints are publicly available and are the core of the OIDC specification
	// =========================================================================
	oidc := rg.Group(router.path)
	{
		// --- OIDC Discovery & Security ---
		oidc.GET("/.well-known/openid-configuration", router.handlers.Discovery.GetConfiguration)
		oidc.GET("/.well-known/jwks.json", router.handlers.JWKS.GetKeys)

		// --- OIDC Core Authorization Flow ---
		oidc.GET("/authorize", router.handlers.OIDC.Authorize)    // Initiates the authorization flow (client frontend request)
		oidc.POST("/token", router.handlers.OIDC.Token)           // Exchange code for a token, update token (client backend request)
		oidc.POST("/revoke", router.handlers.OIDC.Revoke)         // Revocation of a refresh or access token [RFC 7009] (client backend request (Single Session Logout, Client-Wide Logout))
		oidc.POST("/introspect", router.handlers.OIDC.Introspect) // Checking the validity of the token (RFC 7662).
		oidc.GET("/endsession", router.handlers.OIDC.EndSession)  // Initiating a logout from the client side (RP-Initiated Logout) (client frontend request)

		// --- User Interaction & Authentication Endpoints ---
		auth := oidc.Group("/auth")
		{
			auth.POST("/register", router.handlers.User.Register)                               // New user registration (portal frontend request)
			auth.POST("/login/password", router.handlers.User.LoginPassword)                    // Authentication by login/password (sso frontend request with request_id)
			auth.GET("/providers", router.handlers.Provider.List)                               // Getting a list of external providers (Google, etc.).
			auth.GET("/providers/google/callback", router.handlers.Provider.GoogleCallback)     // Callback after authentication with Google
			auth.GET("/providers/telegram/callback", router.handlers.Provider.TelegramCallback) // Callback after authentication with Telegram
		}
	}

	// =========================================================================
	// --- 2. Protected User Portal Endpoints (for AuthFrontPortal) ---
	// Available only to authenticated users via a session cookie
	// Used to manage your profile and sessions
	// =========================================================================
	protected := rg.Group(router.path)
	protected.Use(router.middlewares.UserSession)
	{
		userAPI := protected.Group("/users/me")
		{
			userAPI.GET("", router.handlers.User.GetMyProfile)
			userAPI.GET("/sessions", router.handlers.User.GetMySessions)
			userAPI.DELETE("/sessions/:jti", router.handlers.User.TerminateSession)
			userAPI.POST("/logout-global", router.handlers.User.GlobalLogout)
		}
	}

	// =========================================================================
	// --- 3. Protected Client Endpoints (for AppBack) ---
	// Available to client backends via access_token (Bearer)
	// =========================================================================
	clientProtected := rg.Group(router.path)
	clientProtected.Use(router.middlewares.BearerToken)
	{
		clientProtected.GET("/userinfo", router.handlers.OIDC.UserInfo)
		clientProtected.POST("/logout-client-wide", router.handlers.OIDC.ClientWideLogout) // (Client-Wide Logout)
	}

	// =========================================================================
	// --- 4. Admin API ---
	//Endpoints for system administration. Requires administrator rights
	// =========================================================================
	admin := rg.Group(router.path + "/admin")
	admin.Use(router.middlewares.Admin)
	{
		// --- Client Management (CRUD) ---
		clientsAPI := admin.Group("/clients")
		{
			clientsAPI.GET("", router.handlers.Admin.ListClients)
			clientsAPI.POST("", router.handlers.Admin.CreateClient)
			clientsAPI.GET("/:id", router.handlers.Admin.GetClient)
			clientsAPI.PUT("/:id", router.handlers.Admin.UpdateClient)
			clientsAPI.DELETE("/:id", router.handlers.Admin.DeleteClient)
		}

		// --- User Management (CRUD) ---
		usersAPI := admin.Group("/users")
		{
			usersAPI.GET("", router.handlers.Admin.ListUsers)
			usersAPI.POST("", router.handlers.Admin.CreateUser)
			usersAPI.GET("/:id", router.handlers.Admin.GetUser)
			usersAPI.PUT("/:id", router.handlers.Admin.UpdateUser)
			usersAPI.DELETE("/:id", router.handlers.Admin.DeleteUser)

			// --- Admin-level Session Management ---
			usersAPI.DELETE("/:userID/sessions/:jti", router.handlers.Admin.TerminateUserSession)
			usersAPI.POST("/:userID/logout-global", router.handlers.Admin.GlobalLogoutUser)
		}

		// --- Cryptographic Keys Management ---
		keysAPI := admin.Group("/keys")
		{
			keysAPI.GET("", router.handlers.Admin.GetKeys)            // Get a list of all cryptographic keys (JWKS)
			keysAPI.POST("/rotate", router.handlers.Admin.RotateKeys) // Initiate rotation of token signing keys
		}

		// --- Monitoring ---
		admin.GET("/audit-logs", router.handlers.Admin.GetAuditLogs) // View security audit logs
	}
}
