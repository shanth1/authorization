package telegram

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shanth1/authorization/internal/core/domain"
)

type TelegramProvider struct {
	bot     *tgbotapi.BotAPI
	botName string
}

func NewTelegramProvider(token string, botName string) (*TelegramProvider, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramProvider{bot: bot, botName: botName}, nil
}

func (p *TelegramProvider) Type() string {
	return "telegram"
}

func (p *TelegramProvider) BuildAuthURL(state string) string {
	return fmt.Sprintf(
		"https://t.me/%s?start=%s",
		p.botName,
		state,
	)
}

func (p *TelegramProvider) Exchange(
	ctx context.Context,
	code string, // В Telegram это "start" параметр
) (*domain.User, error) {
	// 1. Получение данных пользователя через Telegram API
	userUpdate, err := p.bot.GetUpdates(tgbotapi.UpdateConfig{
		Timeout: 30,
		Offset:  0,
	})
	if err != nil {
		return nil, err
	}

	// 2. Поиск нужного сообщения по коду
	var userID int64
	for _, update := range userUpdate {
		if update.Message != nil && update.Message.Text == "/start "+code {
			userID = update.Message.From.ID
			break
		}
	}
	if userID == 0 {
		return nil, errors.New("telegram user not found")
	}

	// 3. Формирование объекта пользователя
	return &domain.User{
		Providers: []domain.Provider{
			{
				Type:        "telegram",
				ExternalID:  fmt.Sprint(userID),
				DisplayName: fmt.Sprintf("Telegram User %d", userID),
			},
		},
	}, nil
}
