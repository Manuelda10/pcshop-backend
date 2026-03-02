package http

import (
	"auth-service/config"
	"auth-service/internal/adapters/input/http/dto"
	"strings"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// RegisterRoutes registra todas las rutas del servicio de auth en la app de Fiber.
func RegisterRoutes(app *fiber.App, handler *AuthHandler, cfg config.Config) {
	auth := app.Group("/auth")

	// Rutas públicas — no requieren JWT
	auth.Post("/register", handler.Register)
	auth.Post("/verify-email", handler.VerifyEmail)
	auth.Post("/login", handler.Login)
	auth.Post("/logout", handler.Logout)
	auth.Post("/refresh", handler.RefreshToken)

	// Google OAuth2
	auth.Get("/google", handler.GoogleLogin)
	auth.Get("/google/callback", handler.GoogleCallback)

	// Rutas protegidas — requieren JWT válido
	protected := auth.Group("", JWTMiddleware(cfg.JWT.AccessSecret))
	protected.Get("/me", handler.Me)
}

// JWTMiddleware valida el access token en el header Authorization.
// Si es válido, inyecta el userID en los locals de Fiber para los handlers.
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.FromCtx(c).With(logger.Layer("middleware"), logger.Component("jwt"))

		// Extraer token del header: "Bearer <token>"
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warn("missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error: "missing authorization header",
				Code:  "MISSING_AUTH_HEADER",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			log.Warn("invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error: "invalid authorization header format — expected: Bearer <token>",
				Code:  "INVALID_AUTH_HEADER",
			})
		}

		tokenString := parts[1]

		// Parsear y validar el JWT
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			log.Warn("invalid or expired jwt token", logger.Err(err))
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error: "invalid or expired token",
				Code:  "INVALID_TOKEN",
			})
		}

		// Extraer claims e inyectar en el contexto
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Warn("failed to parse jwt claims")
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error: "invalid token claims",
				Code:  "INVALID_TOKEN_CLAIMS",
			})
		}

		userID, _ := claims["user_id"].(string)
		email, _ := claims["email"].(string)

		c.Locals("userID", userID)
		c.Locals("email", email)

		log.Debug("jwt validated",
			logger.UserID(userID),
			logger.Email(email),
		)

		return c.Next()
	}
}
