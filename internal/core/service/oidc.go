package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

type OIDCService struct {
	clientRepo   ports.OIDCClientRepository
	authCodeRepo ports.AuthCodeRepository
	tokenService ports.TokenService
	sessionCache ports.SessionCache
}

func NewOIDCService(
	clientRepo ports.OIDCClientRepository,
	authCodeRepo ports.AuthCodeRepository,
	tokenService ports.TokenService,
	sessionCache ports.SessionCache,
) *OIDCService {
	return &OIDCService{
		clientRepo:   clientRepo,
		authCodeRepo: authCodeRepo,
		tokenService: tokenService,
		sessionCache: sessionCache,
	}
}

// HandleAuthorize обрабатывает OIDC-запрос авторизации
func (s *OIDCService) HandleAuthorize(
	ctx context.Context,
	req *domain.AuthorizationRequest,
) (string, error) {
	// 1. Валидация клиента
	client, err := s.clientRepo.Get(ctx, req.ClientID)
	if err != nil {
		return "", errors.New("invalid client")
	}

	// 2. Проверка redirect_uri
	validRedirect := false
	for _, uri := range client.RedirectURIs {
		if uri == req.RedirectURI {
			validRedirect = true
			break
		}
	}
	if !validRedirect {
		return "", errors.New("invalid redirect URI")
	}

	// 3. Генерация кода авторизации
	authCode, err := generateSecureCode(16) // Генерация случайного кода
	if err != nil {
		return "", err
	}
	authData := &domain.AuthCodeData{
		ClientID:    req.ClientID,
		RedirectURI: req.RedirectURI,
		Nonce:       req.Nonce,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	if err := s.authCodeRepo.Store(ctx, authCode, authData); err != nil {
		return "", err
	}

	// 4. Формирование redirect URL
	return req.RedirectURI + "?code=" + authCode + "&state=" + req.State, nil
}

// ExchangeCode обменивает код на токены
func (s *OIDCService) ExchangeCode(
	ctx context.Context,
	code string,
) (*domain.Tokens, error) {
	// 1. Получение данных по коду
	authData, err := s.authCodeRepo.Retrieve(ctx, code)
	if err != nil {
		return nil, errors.New("invalid code")
	}

	// 2. Получение клиента
	client, err := s.clientRepo.Get(ctx, authData.ClientID)
	if err != nil {
		return nil, errors.New("invalid client")
	}

	// 3. Генерация токенов
	tokens, err := s.tokenService.GenerateTokens(
		&domain.User{ID: authData.UserID},
		client,
		authData.Nonce,
	)
	if err != nil {
		return nil, err
	}

	// 4. Сохранение сессии (опционально)
	sessionID, err := generateSessionID()
	if err != nil {
		// TODO:
	}

	if err := s.sessionCache.SetTokens(
		ctx,
		sessionID,
		tokens,
		time.Duration(tokens.ExpiresIn)*time.Second,
	); err != nil {
		// Логируем ошибку, но не прерываем процесс
	}

	return tokens, nil
}

// TODO: move to utils:

// GenerateSecureCode генерирует криптографически безопасный код заданной длины
func generateSecureCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateSessionID генерирует криптографически безопасный ID сессии
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
