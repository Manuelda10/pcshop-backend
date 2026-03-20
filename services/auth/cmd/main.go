package main

import (
	"auth/config"
	httpAdapter "auth/internal/adapters/input/http"
	"auth/internal/adapters/input/http/handler"
	authvalidation "auth/internal/adapters/input/http/validation"
	adapterPostgres "auth/internal/adapters/output/postgres"
	"auth/internal/core/services"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Configuración
	cfg := config.Load()

	// 2. Logger
	log := logger.New(logger.Config{
		Level:       cfg.App.LogLevel,
		Env:         cfg.App.Env,
		ServiceName: cfg.App.Name,
	})
	defer log.Sync()

	log.Info("starting auth service",
		logger.Component("main"),
	)

	// 3. Base de datos
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal("failed to connect to database",
			logger.Component("main"),
			logger.Err(err),
		)
	}
	defer db.Close()

	log.Info("database connected", logger.Component("main"), logger.DBTable("postgres"))

	// 4. Adaptadores de salida (output adapters)
	txManager := adapterPostgres.NewTxManager(db)
	userRepo := adapterPostgres.NewUserRepository(db)
	userConsentRepo := adapterPostgres.NewUserConsentRepository(db)
	userConsentStatusRepo := adapterPostgres.NewUserConsentStatusRepository(db)
	tokenRepo := adapterPostgres.NewTokenRepository(db)
	//emailSender := adapterEmail.NewSMTPSender(cfg.Email)

	// 5. Servicio de dominio (core)
	authService := services.NewAuthService(txManager, userRepo, userConsentStatusRepo, userConsentRepo, tokenRepo, cfg)

	validator := authvalidation.NewAuthValidator()

	// 7. Handlers (cada uno recibe solo lo que necesita)
	handlers := httpAdapter.Handlers{
		Register: handler.NewRegisterHandler(authService, validator),
	}

	// 6. Adaptador de entrada (input adapter) — Fiber
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		// Evitar exponer detalles internos en errores de Fiber
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log := logger.FromCtx(c)
			log.Error("unhandled fiber error", logger.Err(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
				"code":  "INTERNAL_ERROR",
			})
		},
	})

	// Middlewares globales
	app.Use(recover.New()) // recuperar de panics sin caer el servicio
	app.Use(logger.FiberMiddleware(log))

	// Rutas
	httpAdapter.RegisterRoutes(app, handlers, cfg)

	// Healthcheck
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": cfg.App.Name,
		})
	})

	// 7. Arranque con graceful shutdown
	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	log.Info("server listening", logger.Component("main"), logger.Path(addr))

	// Canal para señales del SO (Ctrl+C, SIGTERM en Docker/K8s)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Arrancar en goroutine para no bloquear
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatal("server error", logger.Component("main"), logger.Err(err))
		}
	}()

	// Esperar señal de apagado
	<-quit

	log.Info("shutting down gracefully...", logger.Component("main"))

	// Dar 10 segundos para que terminen las requests en vuelo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("forced shutdown", logger.Component("main"), logger.Err(err))
	}

	log.Info("server stopped", logger.Component("main"))
}

// connectDB establece el pool de conexiones a PostgreSQL.
func connectDB(cfg config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		return nil, err
	}

	// Verificar conectividad
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
