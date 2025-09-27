package ports

import (
	"context"

	"github.com/shanth1/authorization/internal/core/domain"
)

type BackchannelSender interface {
	Send(ctx context.Context, client domain.Client, logoutToken string) error
}

type BackchannelReceiver interface {
	Handle(ctx context.Context, logoutToken string) error
}
