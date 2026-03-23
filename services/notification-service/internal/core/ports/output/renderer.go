package output

import "notification-service/internal/core/domain"

// TemplateRenderer es el port de salida para renderizar templates de email.
type TemplateRenderer interface {
	Render(eventType domain.NotificationType, data map[string]string) (subject string, html string, err error)
}
