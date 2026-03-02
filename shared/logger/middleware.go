package logger

import (
	"bytes"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

// FiberMiddleware retorna un middleware de Fiber que loguea automáticamente
// cada request con: method, path, status, latencia, ip, request_id y body (en debug).
//
// Uso en main.go:
//
//	app.Use(logger.FiberMiddleware(log))
func FiberMiddleware(l *Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Asignar o propagar request_id
		requestID := c.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(requestIDHeader, requestID)

		// Logger enriquecido con contexto del request — se usa en todo el ciclo de vida
		reqLogger := l.With(
			RequestID(requestID),
			Method(c.Method()),
			Path(c.Path()),
			IP(c.IP()),
			UserAgent(c.Get("User-Agent")),
		)

		// Almacenar en el contexto de Fiber (para FromCtx en handlers)
		c.Locals("logger", reqLogger)

		// Almacenar también en el context.Context de Go (para FromContext en services y repositories)
		// Así las capas internas no necesitan conocer Fiber en absoluto.
		enrichedCtx := WithContext(c.UserContext(), reqLogger)
		c.SetUserContext(enrichedCtx)

		// Log de entrada — nivel debug para no saturar en producción
		reqLogger.Debug("→ request received",
			Layer("http"),
			zap.String("query_params", string(c.Request().URI().QueryString())),
			zap.ByteString("body", sanitizeBody(c.Body())),
		)

		// Ejecutar el handler
		err := c.Next()

		latency := time.Since(start).Milliseconds()
		statusCode := c.Response().StatusCode()

		// Seleccionar nivel según status code
		switch {
		case statusCode >= 500:
			reqLogger.Error("← response sent",
				Layer("http"),
				StatusCode(statusCode),
				Latency(latency),
				zap.ByteString("response_body", c.Response().Body()),
			)
		case statusCode >= 400:
			reqLogger.Warn("← response sent",
				Layer("http"),
				StatusCode(statusCode),
				Latency(latency),
				zap.ByteString("response_body", c.Response().Body()),
			)
		default:
			reqLogger.Info("← response sent",
				Layer("http"),
				StatusCode(statusCode),
				Latency(latency),
			)
		}

		return err
	}
}

// FromCtx extrae el logger contextualizado desde el contexto de Fiber.
// Si no existe (p.ej. en tests), retorna nil — el caller debe manejarlo.
//
// Uso en handlers:
//
//	log := logger.FromCtx(c)
//	log.Info("processing request", logger.Operation("register_user"))
func FromCtx(c *fiber.Ctx) *Logger {
	l, ok := c.Locals("logger").(*Logger)
	if !ok {
		return nil
	}
	return l
}

// sanitizeBody limpia el body para evitar loguear campos sensibles como passwords.
// En producción podrías ampliar esta lógica con un JSON masker.
func sanitizeBody(body []byte) []byte {
	if len(body) == 0 {
		return []byte("<empty>")
	}

	// Sanitizar campos sensibles básicos
	sensitiveFields := []string{"password", "password_confirmation", "token", "secret"}
	result := string(body)

	for _, field := range sensitiveFields {
		// Reemplaza el valor del campo con *** en JSON
		// Ejemplo: "password":"mipass123" → "password":"***"
		result = replaceSensitiveField(result, field)
	}

	return []byte(result)
}

func replaceSensitiveField(body, field string) string {
	// Busca patrones como "field":"value" o "field": "value"
	patterns := []string{
		`"` + field + `":"`,
		`"` + field + `": "`,
	}

	for _, pattern := range patterns {
		idx := strings.Index(body, pattern)
		if idx == -1 {
			continue
		}
		start := idx + len(pattern)
		end := strings.Index(body[start:], `"`)
		if end == -1 {
			continue
		}
		var buf bytes.Buffer
		buf.WriteString(body[:start])
		buf.WriteString("***")
		buf.WriteString(body[start+end:])
		body = buf.String()
	}

	return body
}
