package service

import (
	"context"
	"fmt"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

// LogoutService manages various logout scenarios
type LogoutService struct {
	sessionService    *SessionService
	tokenRepo         ports.TokenRepository
	clientRepo        ports.ClientRepository
	backchannelSender ports.BackchannelSender
	auditRepo         ports.AuditLogRepository
	jwtSigner         ports.JWTSigner
	random            ports.Random
	clock             ports.Clock
}

func NewLogoutService(
	sessionService *SessionService,
	tokenRepo ports.TokenRepository,
	clientRepo ports.ClientRepository,
	backchannelSender ports.BackchannelSender,
	auditRepo ports.AuditLogRepository,
	jwtSigner ports.JWTSigner,
	random ports.Random,
	clock ports.Clock,
) *LogoutService {
	return &LogoutService{
		sessionService:    sessionService,
		tokenRepo:         tokenRepo,
		clientRepo:        clientRepo,
		backchannelSender: backchannelSender,
		auditRepo:         auditRepo,
		jwtSigner:         jwtSigner,
		random:            random,
		clock:             clock,
	}
}

// LocalLogout - logout from one client (revoke refresh token)
func (s *LogoutService) LocalLogout(
	ctx context.Context,
	refreshTokenJTI domain.JTI,
	clientID domain.ClientID,
) error {
	err := s.tokenRepo.RevokeRefreshByJTI(ctx, refreshTokenJTI)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventTokensRevoked,
		ClientID:  &clientID,
		Timestamp: s.clock.Now(),
		Success:   true,
		Details: map[string]any{
			"type": "local",
		},
	})

	return nil
}

// ClientWideLogout - logout from all sessions of one client
func (s *LogoutService) ClientWideLogout(
	ctx context.Context,
	userID domain.UserID,
	clientID domain.ClientID,
) error {
	tokens, err := s.tokenRepo.ListRefreshByUser(ctx, userID)
	if err != nil {
		return err
	}

	err = s.tokenRepo.RevokeAllByClientUser(ctx, clientID, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	// Send a Back-Channel Logout notification
	client, _ := s.clientRepo.Get(ctx, clientID)
	if client != nil && client.BackchannelLogoutURI != nil {
		for _, token := range tokens {
			if token.ClientID == clientID && !token.Revoked {
				_ = s.sendBackchannelLogout(ctx, client, userID, &token.UserSessionID)
			}
		}
	}

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventTokensRevoked,
		SubjectID: &userID,
		ClientID:  &clientID,
		Timestamp: s.clock.Now(),
		Success:   true,
		Details: map[string]any{
			"type": "client_wide",
		},
	})

	return nil
}

// RemoteSessionTermination - remote session termination from the portal
func (s *LogoutService) RemoteSessionTermination(
	ctx context.Context,
	sessionID domain.UserSessionID,
	userID domain.UserID,
) error {
	session, err := s.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if session.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	tokens, _ := s.tokenRepo.ListRefreshByUserSession(ctx, sessionID)

	err = s.sessionService.TerminateSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Send Back-Channel Logout to all clients
	for _, token := range tokens {
		if !token.Revoked {
			client, _ := s.clientRepo.Get(ctx, token.ClientID)
			if client != nil && client.BackchannelLogoutURI != nil {
				_ = s.sendBackchannelLogout(ctx, client, userID, &sessionID)
			}
		}
	}

	return nil
}

// GlobalLogout - logout from all sessions and all clients
func (s *LogoutService) GlobalLogout(
	ctx context.Context,
	userID domain.UserID,
) error {
	tokens, _ := s.tokenRepo.ListRefreshByUser(ctx, userID)

	// Terminate all sessions (all tokens will be echoed inside)
	err := s.sessionService.GlobalLogout(ctx, userID)
	if err != nil {
		return err
	}

	clientTokens := make(map[domain.ClientID][]domain.UserSessionID)
	for _, token := range tokens {
		if !token.Revoked {
			clientTokens[token.ClientID] = append(clientTokens[token.ClientID], token.UserSessionID)
		}
	}

	// We send a Back-Channel Logout to each client
	for clientID := range clientTokens {
		client, _ := s.clientRepo.Get(ctx, clientID)
		if client != nil && client.BackchannelLogoutURI != nil {
			// Sending without binding to a specific session (all sessions)
			_ = s.sendBackchannelLogout(ctx, client, userID, nil)
		}
	}

	return nil
}

func (s *LogoutService) sendBackchannelLogout(
	ctx context.Context,
	client *domain.Client,
	userID domain.UserID,
	sessionID *domain.UserSessionID,
) error {
	claims := domain.LogoutTokenClaims{
		Issuer:    "https://api.auth.oidc.com", // TODO
		Audience:  string(client.ID),
		IssuedAt:  s.clock.Now(),
		JTI:       domain.JTI(s.random.NewID()),
		Subject:   &userID,
		SessionID: sessionID,
		// TODO:
		Events: map[string]any{
			"http://schemas.openid.net/event/backchannel-logout": map[string]any{},
		},
	}

	logoutToken, err := s.jwtSigner.SignLogoutToken(ctx, claims)
	if err != nil {
		return fmt.Errorf("failed to sign logout token: %w", err)
	}

	err = s.backchannelSender.Send(ctx, *client, logoutToken)
	if err != nil {
		errStr := err.Error()
		_ = s.auditRepo.Log(ctx, &domain.AuditLog{
			ID:        s.random.NewID(),
			EventType: domain.EventBackchannelLogoutSent,
			SubjectID: &userID,
			ClientID:  &client.ID,
			Timestamp: s.clock.Now(),
			Success:   false,
			ErrorMsg:  &errStr,
		})
		return err
	}

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventBackchannelLogoutSent,
		SubjectID: &userID,
		ClientID:  &client.ID,
		Timestamp: s.clock.Now(),
		Success:   true,
	})

	return nil
}
