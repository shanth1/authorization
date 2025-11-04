package domain

import "time"

type AuditEventType string

const (
	EventUserRegistered         AuditEventType = "user_registered"
	EventUserLoggedIn           AuditEventType = "user_logged_in"
	EventUserLoggedOut          AuditEventType = "user_logged_out"
	EventAuthorizationRequested AuditEventType = "authorization_requested"
	EventAuthorizationGranted   AuditEventType = "authorization_granted"
	EventAuthorizationDenied    AuditEventType = "authorization_denied"
	EventTokensIssued           AuditEventType = "tokens_issued"
	EventTokensRevoked          AuditEventType = "tokens_revoked"
	EventTokensRefreshed        AuditEventType = "tokens_refreshed"
	EventConsentGranted         AuditEventType = "consent_granted"
	EventConsentRevoked         AuditEventType = "consent_revoked"
	EventBackchannelLogoutSent  AuditEventType = "backchannel_logout_sent"
	EventSessionTerminated      AuditEventType = "session_terminated"
	EventClientCreated          AuditEventType = "client_created"
	EventClientUpdated          AuditEventType = "client_updated"
	EventClientDeleted          AuditEventType = "client_deleted"
	EventKeysRotated            AuditEventType = "keys_rotated"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string         `json:"id"`
	EventType AuditEventType `json:"event_type"`
	ActorID   *UserID        `json:"actor_id"`   // Who performed the action (nil for system events)
	SubjectID *UserID        `json:"subject_id"` // On whom/what the action was performed
	ClientID  *ClientID      `json:"client_id"`
	IPAddress string         `json:"ip_address"`
	UserAgent string         `json:"user_agent"`
	Details   map[string]any `json:"details"` // Additional event data
	Timestamp time.Time      `json:"timestamp"`
	Success   bool           `json:"success"`   // Was the action completed successfully?
	ErrorMsg  *string        `json:"error_msg"` // Error message if not successful
}
