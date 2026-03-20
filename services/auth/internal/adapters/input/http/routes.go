package http

import (
	"auth/config"
	"auth/internal/adapters/input/http/dto"
	"auth/internal/adapters/input/http/handler"
	"strings"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Handlers struct {
	Register *handler.RegisterHandler
	// Login       *handler.LoginHandler       ← futuro
	// Refresh     *handler.RefreshHandler      ← futuro
}

func RegisterRoutes(app *fiber.App, h Handlers, cfg config.Config) {
	auth := app.Group("/auth")

	auth.Post("/register", h.Register.Handle)

	// OAuth
	//auth.Get("/google/login", h.GoogleOAuth.Login)
	//auth.Get("/google/callback", h.GoogleOAuth.Callback)

	// Futuro:
	// auth.Post("/login", h.Login.Handle)
	// auth.Post("/refresh", h.Refresh.Handle)
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
