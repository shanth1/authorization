package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

// TODO: ?? REMOVE ??

// EventPublisher is responsible for publishing domain events for listeners to consume.
// This is used for auditing, notifications, and other side effects.
type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
