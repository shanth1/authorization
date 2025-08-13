// cmd/telegram_bot/main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	BotToken       string
	AuthServiceURL string
	InternalSecret string
}

func main() {
	log.Println("Starting telegram-bot")
	_ = godotenv.Load()

	cfg := loadConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot API: %v", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		if update.Message.Command() == "start" {
			state := update.Message.CommandArguments()
			if state == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please initiate login from the website first.")
				bot.Send(msg)
				continue
			}

			log.Printf("Received start command with state: %s", state)

			err := completeLogin(cfg, state, update.Message.From.ID, update.Message.From.UserName)
			if err != nil {
				log.Printf("Failed to call auth service: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred during authentication. Please try again.")
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Success! You can now return to your browser.")
			bot.Send(msg)
		}
	}
}

func loadConfig() *Config {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	internalSecret := os.Getenv("INTERNAL_API_SECRET")

	if botToken == "" || authServiceURL == "" || internalSecret == "" {
		log.Fatal("One or more required environment variables are not set: TELEGRAM_BOT_TOKEN, AUTH_SERVICE_URL, INTERNAL_API_SECRET")
	}

	return &Config{
		BotToken:       botToken,
		AuthServiceURL: authServiceURL,
		InternalSecret: internalSecret,
	}
}

func completeLogin(cfg *Config, state string, telegramID int64, username string) error {
	payload := map[string]interface{}{
		"state":       state,
		"telegram_id": telegramID,
		"username":    username,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := cfg.AuthServiceURL + "/api/v1/internal/telegram/complete"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Secret", cfg.InternalSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned non-200 status: %s", resp.Status)
	}

	log.Printf("Successfully completed login for state: %s", state)
	return nil
}
