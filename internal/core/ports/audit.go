package ports

import (
	"context"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type AuditLogFilter struct {
	EventTypes []domain.AuditEventType
	ActorID    *domain.UserID
	SubjectID  *domain.UserID
	ClientID   *domain.ClientID
	StartTime  *time.Time
	EndTime    *time.Time
	Limit      int
	Offset     int
}

// AuditLogRepository manages the audit log
type AuditLogRepository interface {
	Log(ctx context.Context, log *domain.AuditLog) error
	Query(ctx context.Context, filter AuditLogFilter) ([]*domain.AuditLog, error)
}
