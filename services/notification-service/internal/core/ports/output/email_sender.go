package output

import (
	"context"
	"notification-service/internal/core/domain"
)

// EmailSender es el port de salida para enviar emails.
type EmailSender interface {
	Send(ctx context.Context, email domain.Email) error
}
