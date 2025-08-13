package httphandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shanth1/authorization/internal/core/ports"
	"github.com/shanth1/authorization/internal/core/service"
)

type Handler struct {
	authService ports.AuthService
}

func NewHandler(authService ports.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) InitRoutes(router *gin.Engine, internalSecret string) {
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/telegram/initiate", h.telegramInitiate)
			auth.GET("/session/poll", h.pollLoginStatus)
		}

		internal := api.Group("/internal", h.internalAuthMiddleware(internalSecret))
		{
			internal.POST("/telegram/complete", h.telegramComplete)
		}
	}
}

func (h *Handler) telegramInitiate(c *gin.Context) {
	loginURL, err := h.authService.InitiateTelegramLogin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate login"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"login_url": loginURL})
}

func (h *Handler) pollLoginStatus(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State parameter is required"})
		return
	}

	tokens, err := h.authService.PollLoginStatus(c.Request.Context(), state)
	if err != nil {
		if errors.Is(err, service.ErrLoginPending) {
			// 202 Accepted: The request has been accepted for processing,
			// but the processing has not been completed.
			c.JSON(http.StatusAccepted, gin.H{"status": "pending"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to poll status"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// DTO for telegramComplete
type completeLoginRequest struct {
	State      string `json:"state" binding:"required"`
	TelegramID int64  `json:"telegram_id" binding:"required"`
	Username   string `json:"username"`
}

func (h *Handler) telegramComplete(c *gin.Context) {
	var req completeLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.CompleteTelegramLogin(c.Request.Context(), req.State, req.TelegramID, req.Username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()}) // Forbidden is better than Internal error here
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Middleware to protect internal API endpoints
func (h *Handler) internalAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		providedSecret := c.GetHeader("X-Internal-Secret")
		if providedSecret == "" || providedSecret != secret {
			// TODO: logger
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
