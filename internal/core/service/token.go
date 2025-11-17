package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

// TokenService manages the lifecycle of tokens
type TokenService struct {
	codeRepo        ports.AuthorizationCodeRepository
	tokenRepo       ports.TokenRepository
	userSessionRepo ports.UserSessionRepository
	clientRepo      ports.ClientRepository
	auditRepo       ports.AuditLogRepository
	jwtSigner       ports.JWTSigner
	pkceService     ports.PKCEService
	random          ports.Random
	clock           ports.Clock

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewTokenService(
	codeRepo ports.AuthorizationCodeRepository,
	tokenRepo ports.TokenRepository,
	userSessionRepo ports.UserSessionRepository,
	clientRepo ports.ClientRepository,
	auditRepo ports.AuditLogRepository,
	jwtSigner ports.JWTSigner,
	pkceService ports.PKCEService,
	random ports.Random,
	clock ports.Clock,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *TokenService {
	return &TokenService{
		codeRepo:        codeRepo,
		tokenRepo:       tokenRepo,
		userSessionRepo: userSessionRepo,
		clientRepo:      clientRepo,
		auditRepo:       auditRepo,
		jwtSigner:       jwtSigner,
		pkceService:     pkceService,
		random:          random,
		clock:           clock,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// ExchangeCodeForTokens handles POST /token with grant_type=authorization_code
func (s *TokenService) ExchangeCodeForTokens(
	ctx context.Context,
	code domain.Code,
	clientID domain.ClientID,
	clientSecret string,
	redirectURI string,
	codeVerifier *string,
) (*TokenResponse, error) {
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil || !client.Active {
		return nil, domain.ErrInvalidClient
	}

	// TODO: check clientSecret (hash comparison)

	authCode, err := s.codeRepo.Get(ctx, code)
	if err != nil {
		return nil, domain.ErrInvalidGrant
	}

	if authCode.Used {
		return nil, domain.ErrCodeAlreadyUsed
	}

	if s.clock.Now().After(authCode.ExpiresAt) {
		return nil, domain.ErrCodeExpired
	}

	if authCode.ClientID != clientID {
		return nil, domain.ErrInvalidClient
	}

	if authCode.RedirectURI != redirectURI {
		return nil, domain.ErrInvalidRedirectURI
	}

	if authCode.PKCE != nil {
		if codeVerifier == nil {
			return nil, domain.ErrInvalidCodeVerifier
		}
		valid := s.pkceService.Verify(
			*codeVerifier,
			authCode.PKCE.CodeChallenge,
			authCode.PKCE.CodeChallengeMethod,
		)
		if !valid {
			return nil, domain.ErrInvalidCodeVerifier
		}
	}

	err = s.codeRepo.MarkAsUsed(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to mark code as used: %w", err)
	}

	userSession, err := s.userSessionRepo.Get(ctx, authCode.UserSessionID)
	if err != nil {
		return nil, fmt.Errorf("user session not found: %w", err)
	}

	now := s.clock.Now()

	// Access Token
	accessTokenID := domain.TokenID(s.random.NewID())
	accessToken := &domain.AccessToken{
		ID:        accessTokenID,
		UserID:    authCode.UserID,
		ClientID:  authCode.ClientID,
		Scope:     authCode.Scope,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.accessTokenTTL),
		JTI:       domain.JTI(s.random.NewID()),
	}

	// Refresh Token
	refreshTokenID := domain.TokenID(s.random.NewID())
	refreshToken := &domain.RefreshToken{
		ID:            refreshTokenID,
		UserID:        authCode.UserID,
		ClientID:      authCode.ClientID,
		UserSessionID: authCode.UserSessionID,
		Scope:         authCode.Scope,
		IssuedAt:      now,
		ExpiresAt:     now.Add(s.refreshTokenTTL),
		JTI:           domain.JTI(s.random.NewID()),
		Revoked:       false,
	}

	// ID Token (JWT)
	idTokenClaims := domain.IDTokenClaims{
		Issuer:   "https://api.auth.oidc.com",
		Subject:  authCode.UserID,
		Audience: []domain.ClientID{authCode.ClientID},
		Expires:  now.Add(s.accessTokenTTL),
		IssuedAt: now,
		AuthTime: userSession.CreatedAt,
		Nonce:    authCode.Nonce,
		// TODO: add email, name, picture из User
	}

	idToken, err := s.jwtSigner.SignIDToken(ctx, idTokenClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign id_token: %w", err)
	}

	err = s.tokenRepo.SaveAccess(ctx, accessToken, s.accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to save access token: %w", err)
	}

	err = s.tokenRepo.SaveRefresh(ctx, refreshToken, s.refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventTokensIssued,
		SubjectID: &authCode.UserID,
		ClientID:  &authCode.ClientID,
		Timestamp: now,
		Success:   true,
	})

	return &TokenResponse{
		AccessToken:  string(accessToken.JTI), // or JWT itself
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTokenTTL.Seconds()),
		RefreshToken: string(refreshToken.JTI),
		IDToken:      idToken,
	}, nil
}

// RefreshTokens handles POST /token with grant_type=refresh_token
func (s *TokenService) RefreshTokens(
	ctx context.Context,
	refreshTokenValue string,
	clientID domain.ClientID,
	clientSecret string,
) (*TokenResponse, error) {
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil || !client.Active {
		return nil, domain.ErrInvalidClient
	}

	// TODO: check clientSecret

	refreshToken, err := s.tokenRepo.GetRefreshByJTI(ctx, domain.JTI(refreshTokenValue))
	if err != nil {
		return nil, domain.ErrInvalidGrant
	}

	if refreshToken.Revoked {
		return nil, domain.ErrTokenRevoked
	}

	if s.clock.Now().After(refreshToken.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	if refreshToken.ClientID != clientID {
		return nil, domain.ErrInvalidClient
	}

	userSession, err := s.userSessionRepo.Get(ctx, refreshToken.UserSessionID)
	if err != nil || s.clock.Now().After(userSession.ExpiresAt) {
		return nil, fmt.Errorf("user session expired")
	}

	err = s.tokenRepo.RevokeRefreshByJTI(ctx, refreshToken.JTI)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	now := s.clock.Now()

	// New Access Token
	newAccessToken := &domain.AccessToken{
		ID:        domain.TokenID(s.random.NewID()),
		UserID:    refreshToken.UserID,
		ClientID:  refreshToken.ClientID,
		Scope:     refreshToken.Scope,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.accessTokenTTL),
		JTI:       domain.JTI(s.random.NewID()),
	}

	// New Refresh Token
	newRefreshToken := &domain.RefreshToken{
		ID:            domain.TokenID(s.random.NewID()),
		UserID:        refreshToken.UserID,
		ClientID:      refreshToken.ClientID,
		UserSessionID: refreshToken.UserSessionID,
		Scope:         refreshToken.Scope,
		IssuedAt:      now,
		ExpiresAt:     now.Add(s.refreshTokenTTL),
		JTI:           domain.JTI(s.random.NewID()),
		RotatedFrom:   &refreshToken.ID,
		Revoked:       false,
	}

	// New ID Token
	idTokenClaims := domain.IDTokenClaims{
		Issuer:   "https://api.auth.oidc.com",
		Subject:  refreshToken.UserID,
		Audience: []domain.ClientID{refreshToken.ClientID},
		Expires:  now.Add(s.accessTokenTTL),
		IssuedAt: now,
		AuthTime: userSession.CreatedAt,
	}

	idToken, err := s.jwtSigner.SignIDToken(ctx, idTokenClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign id_token: %w", err)
	}

	_ = s.tokenRepo.SaveAccess(ctx, newAccessToken, s.accessTokenTTL)
	_ = s.tokenRepo.SaveRefresh(ctx, newRefreshToken, s.refreshTokenTTL)

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventTokensRefreshed,
		SubjectID: &refreshToken.UserID,
		ClientID:  &refreshToken.ClientID,
		Timestamp: now,
		Success:   true,
	})

	return &TokenResponse{
		AccessToken:  string(newAccessToken.JTI),
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTokenTTL.Seconds()),
		RefreshToken: string(newRefreshToken.JTI),
		IDToken:      idToken,
	}, nil
}

// RevokeToken handles POST /revoke
func (s *TokenService) RevokeToken(
	ctx context.Context,
	tokenValue string,
	clientID domain.ClientID,
	clientSecret string,
) error {
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil || !client.Active {
		return domain.ErrInvalidClient
	}

	// Trying to revoke the refresh token
	jti := domain.JTI(tokenValue)
	err = s.tokenRepo.RevokeRefreshByJTI(ctx, jti)
	if err != nil {
		// Token not found or already revoked - not an error per RFC 7009
		return nil
	}

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventTokensRevoked,
		ClientID:  &clientID,
		Timestamp: s.clock.Now(),
		Success:   true,
	})

	return nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}
