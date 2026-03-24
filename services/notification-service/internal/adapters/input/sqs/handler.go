package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"

	"notification-service/internal/core/domain"
	"notification-service/internal/core/ports/input"

	logger "github.com/Manuelda10/pcshop-backend/shared/logger"
)

// Handler adapta eventos SQS al port de entrada del servicio.
type Handler struct {
	processor input.NotificationProcessor
}

func NewHandler(processor input.NotificationProcessor) *Handler {
	return &Handler{processor: processor}
}

// snsEnvelope representa el wrapper que SNS agrega al mensaje
// cuando lo publica a una cola SQS suscrita.
type snsEnvelope struct {
	Message string `json:"Message"`
}

// HandleBatch procesa un batch de mensajes SQS.
// Retorna BatchItemFailures para que SQS solo reintente los fallidos
// (requiere ReportBatchItemFailures en el event source mapping de la Lambda).
func (h *Handler) HandleBatch(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("handler"),
		logger.Component("sqs"),
		zap.Int("batch_size", len(sqsEvent.Records)),
	)

	log.Info("processing SQS batch")

	var failures []events.SQSBatchItemFailure

	for _, record := range sqsEvent.Records {
		if err := h.processRecord(ctx, log, record); err != nil {
			log.Error("failed to process record",
				logger.Err(err),
				zap.String("message_id", record.MessageId),
			)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}
	}

	if len(failures) > 0 {
		log.Warn("batch completed with failures",
			zap.Int("failed", len(failures)),
			zap.Int("total", len(sqsEvent.Records)),
		)
	} else {
		log.Info("batch completed successfully")
	}

	return events.SQSEventResponse{
		BatchItemFailures: failures,
	}, nil
}

func (h *Handler) processRecord(ctx context.Context, log *logger.Logger, record events.SQSMessage) error {
	// Los mensajes de SNS vienen envueltos en un envelope.
	// Intentamos desempaquetar; si no es SNS, usamos el body directo.
	body := record.Body

	var envelope snsEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err == nil && envelope.Message != "" {
		body = envelope.Message
	}

	var event domain.NotificationEvent
	if err := json.Unmarshal([]byte(body), &event); err != nil {
		return fmt.Errorf("unmarshal event: %w", err)
	}

	log.Debug("parsed notification event",
		zap.String("event_type", string(event.Type)),
		zap.String("message_id", record.MessageId),
	)

	return h.processor.Process(ctx, event)
}
