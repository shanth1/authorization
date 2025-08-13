// internal/core/service/auth.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
	"github.com/shanth1/gotools/log"
)

const (
	loginStatePrefix = "login_state:"
	sessionPrefix    = "session:"
	stateTTL         = 5 * time.Minute
)

type AuthService struct {
	userRepo     ports.UserRepository
	stateCache   ports.StateCache
	sessionCache ports.SessionCache
	tokenSvc     ports.TokenService
	botName      string
}

func NewAuthService(userRepo ports.UserRepository, stateCache ports.StateCache, sessionCache ports.SessionCache, tokenSvc ports.TokenService, botName string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		stateCache:   stateCache,
		sessionCache: sessionCache,
		tokenSvc:     tokenSvc,
		botName:      botName,
	}
}

func (s *AuthService) InitiateTelegramLogin(ctx context.Context) (string, error) {
	logger := log.FromContext(ctx)

	state := uuid.NewString()
	key := loginStatePrefix + state

	if err := s.stateCache.SetState(ctx, key, stateTTL); err != nil {
		return "", fmt.Errorf("internal error on state generation")
	}

	loginURL := fmt.Sprintf("https://t.me/%s?start=%s", s.botName, state)
	logger.Info().Msg("Generated telegram login URL")

	return loginURL, nil
}

func (s *AuthService) CompleteTelegramLogin(ctx context.Context, state string, telegramID int64, username string) error {
	logger := log.FromContext(ctx)

	stateKey := loginStatePrefix + state

	if err := s.stateCache.CheckState(ctx, stateKey); err != nil {
		return fmt.Errorf("invalid or expired state")
	}

	user, err := s.userRepo.FindByTelegramID(ctx, telegramID)
	if err != nil {
		logger.Info().Int("telegramID", int(telegramID)).Msg("User not found, creating a new one")
		newUser := &domain.User{TelegramID: telegramID, Username: username}
		if errCreate := s.userRepo.Create(ctx, newUser); errCreate != nil {
			return fmt.Errorf("failed to create user: %w", errCreate)
		}
		user = newUser
	}

	tokens, err := s.tokenSvc.GenerateTokens(user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate tokens: %w", err)
	}

	sessionKey := sessionPrefix + state
	if err := s.sessionCache.SetSession(ctx, sessionKey, tokens, stateTTL); err != nil {
		return fmt.Errorf("set session: %w", err)
	}

	_ = s.stateCache.DeleteState(ctx, stateKey) // Ignore error, it will expire anyway
	logger.Info().Str("userID", user.ID).Msg("Telegram login completed successfully, session is ready for polling")

	return nil
}

var ErrLoginPending = errors.New("login is pending")

func (s *AuthService) PollLoginStatus(ctx context.Context, state string) (*domain.Tokens, error) {
	logger := log.FromContext(ctx)

	sessionKey := sessionPrefix + state

	tokens, err := s.sessionCache.GetSession(ctx, sessionKey)
	if err != nil {
		return nil, ErrLoginPending // Not an error, just not ready yet
	}

	_ = s.sessionCache.DeleteSession(ctx, sessionKey)

	logger.Info().Str("state", state).Msg("Session polled successfully")

	return tokens, nil
}
