package input

import (
	"context"
	"notification-service/internal/core/domain"
)

type NotificationProcessor interface {
	Process(ctx context.Context, event domain.NotificationEvent) error
}
