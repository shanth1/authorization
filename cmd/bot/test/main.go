package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			state := update.Message.CommandArguments()
			if state == "" {
				msg.Text = "Invalid request"
				bot.Send(msg)
				continue
			}

			// Подтверждение входа
			msg.Text = "Подтвердите вход в систему"
			confirmBtn := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "confirm_"+state),
				),
			)
			msg.ReplyMarkup = confirmBtn
			bot.Send(msg)
		}
	}
}

func handleConfirmation(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	data := strings.TrimPrefix(callbackQuery.Data, "confirm_")
	state := data

	// Формируем данные для IDP
	_ = map[string]interface{}{
		"state":      state,
		"id":         callbackQuery.From.ID,
		"first_name": callbackQuery.From.FirstName,
		"last_name":  callbackQuery.From.LastName,
		"username":   callbackQuery.From.UserName,
		"auth_date":  time.Now().Unix(),
	}

	// Отправляем данные на IDP
	resp, err := http.Post("https://your-idp.com/telegram/callback", "application/json", bytes.NewBuffer(nil))
	if err != nil {
		// Обработка ошибки
		return
	}
	defer resp.Body.Close()

	var result struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Обработка ошибки
		return
	}

	// Перенаправляем пользователя
	redirectURL := fmt.Sprintf("https://client-app.com/callback?code=%s", result.Code)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Вы успешно авторизованы! [Перейти в приложение]("+redirectURL+")")
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
