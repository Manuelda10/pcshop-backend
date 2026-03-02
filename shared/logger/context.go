package logger

import (
	"context"

	"go.uber.org/zap"
)

// contextKey es un tipo privado para evitar colisiones con otras librerías
// que también guarden valores en el context.
type contextKey struct{}

var loggerKey = contextKey{}

// WithContext retorna un nuevo contexto con el logger embebido.
// Llamar en el middleware de Fiber para que todas las capas
// puedan acceder al logger enriquecido con el request_id.
//
// Uso en middleware:
//
//	enrichedCtx := logger.WithContext(c.UserContext(), reqLogger)
//	c.SetUserContext(enrichedCtx)
func WithContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext extrae el logger del context.
// Si no existe (tests unitarios, jobs, etc.) retorna un logger funcional
// con nivel info para que el programa no explote silenciosamente.
//
// Uso en service:
//
//	func (s *AuthService) Register(ctx context.Context, cmd RegisterCommand) (*User, error) {
//	    log := logger.FromContext(ctx)
//	    log.Info("registering user", logger.Layer("service"))
//	}
//
// Uso en repository:
//
//	func (r *UserRepo) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
//	    log := logger.FromContext(ctx)
//	    log.Debug("inserting user", logger.Layer("repository"), logger.DBTable("users"))
//	}
func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(loggerKey).(*Logger)
	if !ok || l == nil {
		return newFallbackLogger()
	}
	return l
}

// newFallbackLogger retorna un logger funcional de emergencia.
// Se usa cuando no hay logger en el context — evita nil panics.
// Loguea un warning para que sepas que algo no está propagando el contexto bien.
func newFallbackLogger() *Logger {
	fallback := New(Config{
		Level:       "info",
		Env:         "development",
		ServiceName: "unknown",
	})
	fallback.Warn("logger not found in context, using fallback — check middleware setup",
		zap.String("hint", "ensure logger.FiberMiddleware is registered and ctx is propagated"),
	)
	return fallback
}
