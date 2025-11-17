package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

// SessionService manages SSO user sessions
type SessionService struct {
	userSessionRepo ports.UserSessionRepository
	tokenRepo       ports.TokenRepository
	auditRepo       ports.AuditLogRepository
	random          ports.Random
	clock           ports.Clock
	sessionTTL      time.Duration
}

func NewSessionService(
	userSessionRepo ports.UserSessionRepository,
	tokenRepo ports.TokenRepository,
	auditRepo ports.AuditLogRepository,
	random ports.Random,
	clock ports.Clock,
	sessionTTL time.Duration,
) *SessionService {
	return &SessionService{
		userSessionRepo: userSessionRepo,
		tokenRepo:       tokenRepo,
		auditRepo:       auditRepo,
		random:          random,
		clock:           clock,
		sessionTTL:      sessionTTL,
	}
}

// CreateSession creates a new SSO session after successful authentication
func (s *SessionService) CreateSession(
	ctx context.Context,
	userID domain.UserID,
	provider domain.Provider,
	ipAddress string,
	userAgent string,
) (*domain.UserSession, error) {

	sessionID := domain.UserSessionID(s.random.NewID())
	now := s.clock.Now()

	session := &domain.UserSession{
		ID:         sessionID,
		UserID:     userID,
		Provider:   provider,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  now,
		ExpiresAt:  now.Add(s.sessionTTL),
		LastUsedAt: now,
	}

	err := s.userSessionRepo.Create(ctx, session, s.sessionTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventUserLoggedIn,
		SubjectID: &userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Timestamp: now,
		Success:   true,
		Details: map[string]any{
			"provider": provider,
		},
	})

	return session, nil
}

// GetSession gets the session and updates the last used time
func (s *SessionService) GetSession(
	ctx context.Context,
	sessionID domain.UserSessionID,
) (*domain.UserSession, error) {

	session, err := s.userSessionRepo.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if s.clock.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	_ = s.userSessionRepo.UpdateLastUsed(ctx, sessionID)

	return session, nil
}

// TerminateSession terminates a specific session
func (s *SessionService) TerminateSession(
	ctx context.Context,
	sessionID domain.UserSessionID,
) error {
	session, err := s.userSessionRepo.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	err = s.userSessionRepo.Delete(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	_ = s.tokenRepo.RevokeAllByUserSession(ctx, sessionID)

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventSessionTerminated,
		SubjectID: &session.UserID,
		Timestamp: s.clock.Now(),
		Success:   true,
	})

	return nil
}

// ListUserSessions returns a list of all active user sessions.
func (s *SessionService) ListUserSessions(
	ctx context.Context,
	userID domain.UserID,
) ([]*domain.UserSession, error) {

	return s.userSessionRepo.ListByUser(ctx, userID)
}

// GlobalLogout terminates all user sessions
func (s *SessionService) GlobalLogout(
	ctx context.Context,
	userID domain.UserID,
) error {
	err := s.userSessionRepo.DeleteAllByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	_ = s.tokenRepo.RevokeAllByUser(ctx, userID)

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventUserLoggedOut,
		SubjectID: &userID,
		Timestamp: s.clock.Now(),
		Success:   true,
		Details: map[string]any{
			"type": "global",
		},
	})

	return nil
}
