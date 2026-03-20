package handler

import (
	"auth/internal/core/domain"
	"errors"

	sharedhttp "github.com/Manuelda10/pcshop-backend/shared/http"
	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/gofiber/fiber/v2"
)

func handleServiceError(c *fiber.Ctx, err error) error {
	log := logger.FromCtx(c)

	// Ajusta estos cases a los errores reales de tu dominio
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return sharedhttp.Conflict(c, "EMAIL_EXISTS", "El email ya está registrado")

	case errors.Is(err, domain.ErrDocumentAlreadyExists):
		return sharedhttp.Conflict(c, "DOCUMENT_EXISTS", "El documento ya está registrado")

	case errors.Is(err, domain.ErrInvalidCredentials):
		return sharedhttp.Unauthorized(c, "INVALID_CREDENTIALS", "Credenciales inválidas")

	default:
		log.Error("unhandled service error", logger.Err(err))
		return sharedhttp.InternalError(c)
	}
}
