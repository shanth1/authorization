package httphandler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shanth1/authorization/internal/adapters/input/telegram"
	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

type OIDCHandlers struct {
	oidcService     ports.OIDCService
	tokenService    ports.TokenService
	userRepo        ports.UserRepository
	telegramHandler *telegram.TelegramHandler
	jwks            interface{} // Заглушка для JWKS
}

func NewOIDCHandlers(
	oidcService ports.OIDCService,
	tokenService ports.TokenService,
	userRepo ports.UserRepository,
	telegramHandler *telegram.TelegramHandler,
	jwks interface{},
) *OIDCHandlers {
	return &OIDCHandlers{
		oidcService:     oidcService,
		tokenService:    tokenService,
		userRepo:        userRepo,
		telegramHandler: telegramHandler,
		jwks:            jwks,
	}
}

// AuthorizeHandler обрабатывает OIDC Authorization Request
func (h *OIDCHandlers) AuthorizeHandler(c *gin.Context) {
	var req domain.AuthorizationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	authURL, err := h.oidcService.HandleAuthorize(c.Request.Context(), &req)
	if err != nil {
		handleOIDCError(c, err)
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

// TokenHandler обрабатывает OIDC Token Request
func (h *OIDCHandlers) TokenHandler(c *gin.Context) {
	var req domain.TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	tokens, err := h.oidcService.ExchangeCode(c.Request.Context(), req.Code)
	if err != nil {
		handleOIDCError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// UserInfoHandler возвращает информацию о пользователе
func (h *OIDCHandlers) UserInfoHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := h.tokenService.ParseAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
		return
	}

	// Преобразуем в стандартные OIDC claims
	userInfo := domain.UserInfo{
		Sub:           user.ID,
		Email:         user.PrimaryEmail,
		EmailVerified: true,
	}

	if len(user.Providers) > 0 {
		userInfo.PreferredName = user.Providers[0].DisplayName
		userInfo.Picture = user.Providers[0].AvatarURL
	}

	c.JSON(http.StatusOK, userInfo)
}

// DiscoveryHandler возвращает OIDC Discovery Document
func (h *OIDCHandlers) DiscoveryHandler(c *gin.Context) {
	baseURL := "https://" + c.Request.Host
	c.JSON(http.StatusOK, gin.H{
		"issuer":                                baseURL,
		"authorization_endpoint":                baseURL + "/authorize",
		"token_endpoint":                        baseURL + "/token",
		"userinfo_endpoint":                     baseURL + "/userinfo",
		"jwks_uri":                              baseURL + "/.well-known/jwks.json",
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	})
}

// JWKSHandler возвращает JSON Web Key Set
func (h *OIDCHandlers) JWKSHandler(c *gin.Context) {
	c.JSON(http.StatusOK, h.jwks)
}

// RegisterRoutes регистрирует все OIDC роуты
func (h *OIDCHandlers) RegisterRoutes(router *gin.Engine) {
	router.GET("/authorize", h.AuthorizeHandler)
	router.POST("/token", h.TokenHandler)
	router.GET("/userinfo", h.UserInfoHandler)
	router.GET("/.well-known/openid-configuration", h.DiscoveryHandler)
	router.GET("/.well-known/jwks.json", h.JWKSHandler)

	// Telegram-specific роуты
	router.GET("/auth/telegram/init", h.telegramHandler.InitAuthHandler)
	router.POST("/telegram/webhook", h.telegramHandler.WebhookHandler)
}

func handleOIDCError(c *gin.Context, err error) {
	// ... (реализация обработки ошибок)
}
