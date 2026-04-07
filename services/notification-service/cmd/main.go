package main

import (
	"context"
	"log"
	"notification-service/config"
	"notification-service/internal/adapters/output/ses"
	"notification-service/internal/adapters/output/templates"
	"notification-service/internal/core/services"

	sqsAdapter "notification-service/internal/adapters/input/sqs"

	sharedlogger "github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	// 1. Config
	cfg := config.Load()

	// 2. Logger
	logger := sharedlogger.New(sharedlogger.Config{
		Level:       cfg.App.LogLevel,
		Env:         cfg.App.Env,
		ServiceName: "notification-service",
	})
	defer logger.Sync()

	logger.Info("initializing notification service")

	// 3. AWS SDK config (la Lambda hereda el IAM role automáticamente)
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	// 4. Output adapters
	renderer, err := templates.NewRenderer(cfg.App.Name)
	if err != nil {
		log.Fatalf("failed to compile templates: %v", err)
	}

	sender := ses.NewSender(awsCfg, cfg.Email.FromAddress)

	// 5. Core service
	notificationService := services.NewNotificationService(renderer, sender)

	// 6. Input adapter (SQS handler)
	sqsHandler := sqsAdapter.NewHandler(notificationService)

	// 7. Iniciar Lambda
	// Wrapeamos para inyectar el logger en el context antes de cada invocación.
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
		ctx = sharedlogger.WithContext(ctx, logger)
		return sqsHandler.HandleBatch(ctx, sqsEvent)
	})
}
