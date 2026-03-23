package services

import (
	"context"
	"fmt"
	"notification-service/internal/core/domain"
	"notification-service/internal/core/ports/output"

	logger "github.com/Manuelda10/pcshop-backend/shared/logger"
	"go.uber.org/zap"
)

// requiredFields define qué campos del payload son obligatorios por tipo de evento.
var requiredFields = map[domain.NotificationType][]string{
	domain.NotificationUserRegistered: {
		"email", "firstName", "otpCode",
	},
	domain.NotificationPasswordResetRequest: {
		"email", "firstName", "otpCode",
	},
}

// NotificationService implementa la lógica de negocio para procesar notificaciones.
type NotificationService struct {
	renderer output.TemplateRenderer
	sender   output.EmailSender
}

func NewNotificationService(renderer output.TemplateRenderer, sender output.EmailSender) *NotificationService {
	return &NotificationService{
		renderer: renderer,
		sender:   sender,
	}
}

func (s *NotificationService) Process(ctx context.Context, event domain.NotificationEvent) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("process_notification"),
		zap.String("event_type", string(event.Type)),
	)

	log.Info("processing notification event",
		logger.Email(event.Payload["email"]),
	)

	// 1. Validar tipo de evento conocido
	fields, ok := requiredFields[event.Type]
	if !ok {
		log.Warn("unknown event type")
		return fmt.Errorf("%w: %s", domain.ErrUnknownEventType, event.Type)
	}

	// 2. Validar campos requeridos en el payload
	for _, field := range fields {
		if event.Payload[field] == "" {
			log.Warn("missing payload field", zap.String("field", field))
			return fmt.Errorf("%w: %s for event %s", domain.ErrMissingPayloadField, field, event.Type)
		}
	}

	// 3. Renderizar template
	subject, html, err := s.renderer.Render(event.Type, event.Payload)
	if err != nil {
		log.Error("template render failed", logger.Err(err))
		return fmt.Errorf("%w: %v", domain.ErrTemplateRender, err)
	}

	// 4. Enviar email
	email := domain.Email{
		To:      event.Payload["email"],
		Subject: subject,
		HTML:    html,
	}

	if err := s.sender.Send(ctx, email); err != nil {
		log.Error("email send failed", logger.Err(err), zap.String("to", email.To))
		return fmt.Errorf("%w: %v", domain.ErrEmailSendFailed, err)
	}

	log.Info("notification sent successfully", zap.String("to", email.To))
	return nil
}
