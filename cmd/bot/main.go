// cmd/telegram_bot/main.go
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	botcfg "github.com/shanth1/authorization/internal/config/bot"
	"github.com/shanth1/gotools/ctx"
	"github.com/shanth1/gotools/env"
	"github.com/shanth1/gotools/flags"
	"github.com/shanth1/gotools/log"
)

type startCfg struct {
	EnvPath string `flag:"env-path" usage:"[OPTIONAL] Path to env file"`
}

func main() {
	logger := log.New(log.WithService("auth"))

	ctx, cancel := ctx.GetAppCtx()
	defer cancel()

	var startCfg startCfg
	if err := flags.RegisterFromStruct(&startCfg); err != nil {
		logger.Fatal().Err(err).Msg("Register flags from struct")
	}
	flag.Parse()

	cfg := &botcfg.Config{}
	if err := env.LoadIntoStruct(startCfg.EnvPath, cfg); err != nil {
		logger.Fatal().Err(err).Msg("Load env into struct")
	}

	ctx = log.NewContext(ctx, logger)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create bot API")
	}

	bot.Debug = false
	logger.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
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

				logger.Info().Msgf("Received start command with state: %s", state)

				err := completeLogin(cfg, state, update.Message.From.ID, update.Message.From.UserName)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to call auth service")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred during authentication. Please try again.")
					bot.Send(msg)
					continue
				}

				logger.Info().Msgf("Successfully completed login for state: %s", state)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Success! You can now return to your browser.")
				bot.Send(msg)
			}
		}
	}
}

func completeLogin(cfg *botcfg.Config, state string, telegramID int64, username string) error {
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
	req.Header.Set("X-Internal-Secret", cfg.APIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned non-200 status: %s", resp.Status)
	}

	return nil
}
