package telegram

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
	"github.com/shanth1/authorization/internal/utils"
)

type TelegramHandler struct {
	bot         *tgbotapi.BotAPI
	oidcService ports.OIDCService
	authService ports.AuthService
	stateCache  ports.StateCache
}

func NewTelegramHandler(
	bot *tgbotapi.BotAPI,
	oidcService ports.OIDCService,
	authService ports.AuthService,
	stateCache ports.StateCache,
) *TelegramHandler {
	return &TelegramHandler{
		bot:         bot,
		oidcService: oidcService,
		authService: authService,
		stateCache:  stateCache,
	}
}

// InitAuthHandler инициирует аутентификацию через Telegram
func (h *TelegramHandler) InitAuthHandler(c *gin.Context) {
	// Генерация state-параметра
	state, err := utils.GenerateSecureCode(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	// Сохранение state в кэш
	if err := h.stateCache.Store(c, state, &domain.AuthorizationRequest{
		// Здесь должны быть параметры из запроса
	}, 10*time.Minute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store state"})
		return
	}

	// Формирование ссылки для Telegram-бота
	authURL := fmt.Sprintf("https://t.me/%s?start=%s", h.bot.Self.UserName, state)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
	})
}

// WebhookHandler обрабатывает входящие обновления от Telegram
func (h *TelegramHandler) WebhookHandler(c *gin.Context) {
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Обработка только команд /start
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "start" {
		h.handleStartCommand(c, update.Message)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ignored"})
}

func (h *TelegramHandler) handleStartCommand(c *gin.Context, msg *tgbotapi.Message) {
	// Извлекаем state из аргументов команды
	state := msg.CommandArguments()
	if state == "" {
		log.Println("Missing state in /start command")
		return
	}

	// Получаем данные из кэша
	authReq, err := h.stateCache.Verify(c, state)
	if err != nil {
		log.Printf("Invalid state: %v", err)
		return
	}

	// Создаем фейковый код для завершения аутентификации
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		log.Println("generate session id:", err)
		return
	}
	authCode := "telegram_" + sessionID

	// Сохраняем данные для обмена кода
	authData := &domain.AuthCodeData{
		ClientID:    authReq.ClientID,
		RedirectURI: authReq.RedirectURI,
		UserID:      fmt.Sprint(msg.From.ID), // Используем ID пользователя в Telegram
		Nonce:       authReq.Nonce,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	if err := h.oidcService.(ports.AuthCodeRepository).Store(c, authCode, authData); err != nil {
		log.Printf("Failed to store auth code: %v", err)
		return
	}

	// Отправляем сообщение пользователю
	response := "✅ Вы успешно аутентифицированы! Возвращайтесь в приложение."
	reply := tgbotapi.NewMessage(msg.Chat.ID, response)
	if _, err := h.bot.Send(reply); err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	// Возвращаемся в OIDC flow
	c.Redirect(http.StatusFound, authReq.RedirectURI+"?code="+authCode+"&state="+authReq.State)
}
