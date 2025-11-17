package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/authorization/internal/core/ports"
)

// AuthorizationService manages the authorization flow
type AuthorizationService struct {
	requestRepo     ports.AuthorizationRequestRepository
	userSessionRepo ports.UserSessionRepository
	clientRepo      ports.ClientRepository
	codeRepo        ports.AuthorizationCodeRepository
	consentRepo     ports.ConsentRepository
	auditRepo       ports.AuditLogRepository
	random          ports.Random
	clock           ports.Clock
}

func NewAuthorizationService(
	requestRepo ports.AuthorizationRequestRepository,
	userSessionRepo ports.UserSessionRepository,
	clientRepo ports.ClientRepository,
	codeRepo ports.AuthorizationCodeRepository,
	consentRepo ports.ConsentRepository,
	auditRepo ports.AuditLogRepository,
	random ports.Random,
	clock ports.Clock,
) *AuthorizationService {
	return &AuthorizationService{
		requestRepo:     requestRepo,
		userSessionRepo: userSessionRepo,
		clientRepo:      clientRepo,
		codeRepo:        codeRepo,
		consentRepo:     consentRepo,
		auditRepo:       auditRepo,
		random:          random,
		clock:           clock,
	}
}

// InitiateAuthorizationRequest handles GET /authorize
// Returns request_id and redirect URL
func (s *AuthorizationService) InitiateAuthorizationRequest(
	ctx context.Context,
	clientID domain.ClientID,
	redirectURI string,
	scope []domain.Scope,
	state domain.State,
	nonce *domain.Nonce,
	pkce *domain.PKCE,
	prompt *domain.Prompt,
	maxAge *int,
	userSessionCookie *domain.UserSessionID, // из cookie user_session_id
) (*domain.AuthorizationRequest, error) {
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}
	if !client.Active {
		return nil, domain.ErrInvalidClient
	}

	if !s.isValidRedirectURI(client, redirectURI) {
		return nil, domain.ErrInvalidRedirectURI
	}

	if !s.hasValidScopes(client, scope) {
		return nil, domain.ErrInvalidScope
	}

	requestID := domain.RequestID(s.random.NewID())
	authRequest := &domain.AuthorizationRequest{
		ID:          requestID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		State:       state,
		Nonce:       nonce,
		PKCE:        pkce,
		Prompt:      prompt,
		MaxAge:      maxAge,
		CreatedAt:   s.clock.Now(),
		ExpiresAt:   s.clock.Now().Add(10 * time.Minute),
	}

	if userSessionCookie != nil {
		userSession, err := s.userSessionRepo.Get(ctx, *userSessionCookie)
		if err == nil && s.clock.Now().Before(userSession.ExpiresAt) {
			authRequest.UserSessionID = &userSession.ID

			_ = s.userSessionRepo.UpdateLastUsed(ctx, userSession.ID)

			needsReauth := s.needsReauthentication(prompt, maxAge, userSession)
			if !needsReauth {
				// The user is already logged in and does not require re-authentication
				// Checking consent
				needsConsent := s.needsConsent(ctx, client, userSession.UserID, scope)
				if !needsConsent {
					authRequest.ConsentGranted = true
				}
			} else {
				// Re-authentication required
				authRequest.UserSessionID = nil
			}
		}
	}

	err = s.requestRepo.Save(ctx, authRequest, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to save authorization request: %w", err)
	}

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventAuthorizationRequested,
		ClientID:  &clientID,
		Timestamp: s.clock.Now(),
		Success:   true,
	})

	return authRequest, nil
}

// CompleteAuthentication is called after the user has been successfully authenticated.
func (s *AuthorizationService) CompleteAuthentication(
	ctx context.Context,
	requestID domain.RequestID,
	userSessionID domain.UserSessionID,
) error {
	authRequest, err := s.requestRepo.Get(ctx, requestID)
	if err != nil {
		return fmt.Errorf("authorization request not found: %w", err)
	}

	authRequest.UserSessionID = &userSessionID

	return s.requestRepo.Update(ctx, authRequest)
}

// CompleteConsent is called after the user consents.
func (s *AuthorizationService) CompleteConsent(
	ctx context.Context,
	requestID domain.RequestID,
	granted bool,
	scopes []domain.Scope,
) error {
	authRequest, err := s.requestRepo.Get(ctx, requestID)
	if err != nil {
		return fmt.Errorf("authorization request not found: %w", err)
	}

	if !granted {
		// The user denied access
		_ = s.auditRepo.Log(ctx, &domain.AuditLog{
			ID:        s.random.NewID(),
			EventType: domain.EventAuthorizationDenied,
			ClientID:  &authRequest.ClientID,
			Timestamp: s.clock.Now(),
			Success:   true,
		})
		return fmt.Errorf("user denied consent")
	}

	authRequest.ConsentGranted = true

	if authRequest.UserSessionID != nil {
		userSession, _ := s.userSessionRepo.Get(ctx, *authRequest.UserSessionID)
		if userSession != nil {
			consent := &domain.ConsentGrant{
				ID:        s.random.NewID(),
				UserID:    userSession.UserID,
				ClientID:  authRequest.ClientID,
				Scopes:    scopes,
				GrantedAt: s.clock.Now(),
			}
			_ = s.consentRepo.Save(ctx, consent)
		}
	}

	return s.requestRepo.Update(ctx, authRequest)
}

// IssueAuthorizationCode выдает код авторизации
func (s *AuthorizationService) IssueAuthorizationCode(
	ctx context.Context,
	requestID domain.RequestID,
) (*domain.AuthorizationCode, error) {
	authRequest, err := s.requestRepo.Get(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("authorization request not found: %w", err)
	}

	if authRequest.UserSessionID == nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	userSession, err := s.userSessionRepo.Get(ctx, *authRequest.UserSessionID)
	if err != nil {
		return nil, fmt.Errorf("user session not found: %w", err)
	}

	client, _ := s.clientRepo.Get(ctx, authRequest.ClientID)
	if client.RequiresConsent && !authRequest.ConsentGranted {
		return nil, fmt.Errorf("consent not granted")
	}

	codeValue := domain.Code(s.random.NewID())
	code := &domain.AuthorizationCode{
		Code:          codeValue,
		RequestID:     requestID,
		UserSessionID: *authRequest.UserSessionID,
		UserID:        userSession.UserID,
		ClientID:      authRequest.ClientID,
		RedirectURI:   authRequest.RedirectURI,
		Scope:         authRequest.Scope,
		Nonce:         authRequest.Nonce,
		PKCE:          authRequest.PKCE,
		IssuedAt:      s.clock.Now(),
		ExpiresAt:     s.clock.Now().Add(60 * time.Second),
		Used:          false,
	}

	err = s.codeRepo.Save(ctx, code, 60*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to save authorization code: %w", err)
	}

	_ = s.requestRepo.Delete(ctx, requestID)

	_ = s.auditRepo.Log(ctx, &domain.AuditLog{
		ID:        s.random.NewID(),
		EventType: domain.EventAuthorizationGranted,
		SubjectID: &userSession.UserID,
		ClientID:  &authRequest.ClientID,
		Timestamp: s.clock.Now(),
		Success:   true,
	})

	return code, nil
}

func (s *AuthorizationService) isValidRedirectURI(client *domain.Client, uri string) bool {
	for _, allowed := range client.RedirectURIs {
		if string(allowed) == uri {
			return true
		}
	}
	return false
}

func (s *AuthorizationService) hasValidScopes(client *domain.Client, requested []domain.Scope) bool {
	for _, req := range requested {
		found := false
		for _, allowed := range client.Scopes {
			if req == allowed {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *AuthorizationService) needsReauthentication(
	prompt *domain.Prompt,
	maxAge *int,
	userSession *domain.UserSession,
) bool {
	// If prompt=login, re-authentication is required
	if prompt != nil && *prompt == domain.PromptLogin {
		return true
	}

	if maxAge != nil {
		authAge := time.Since(userSession.CreatedAt).Seconds()
		if authAge > float64(*maxAge) {
			return true
		}
	}

	return false
}

func (s *AuthorizationService) needsConsent(
	ctx context.Context,
	client *domain.Client,
	userID domain.UserID,
	scopes []domain.Scope,
) bool {
	// First-party clients do not require consent
	if client.Type == domain.ClientTypeFirstParty {
		return false
	}

	if !client.RequiresConsent {
		return false
	}

	// If prompt=consent, always show
	// (this is checked in the calling code)

	hasConsent, _ := s.consentRepo.Exists(ctx, userID, client.ID, scopes)
	return !hasConsent
}
