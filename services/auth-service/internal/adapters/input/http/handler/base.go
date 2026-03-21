package handler

import (
	sharedhttp "github.com/Manuelda10/pcshop-backend/shared/http"
	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/Manuelda10/pcshop-backend/shared/validation"
	"github.com/gofiber/fiber/v2"
)

// Base contiene las dependencias compartidas por todos los handlers de auth.
type Base struct {
	Validator *validation.Validator
}

// parseAndValidate es un helper genérico: parsea body + valida.
// Retorna true si hubo error (y ya respondió al cliente).
func parseAndValidate[T any](c *fiber.Ctx, b *Base, req *T) bool {
	log := logger.FromCtx(c)

	if err := c.BodyParser(req); err != nil {
		log.Warn("invalid request body", logger.Err(err))
		sharedhttp.BadRequest(c, sharedhttp.CodeInvalidBody, "Cuerpo de la solicitud inválido")
		return true
	}

	if fieldErrors := b.Validator.Validate(*req); fieldErrors != nil {
		sharedhttp.ValidationError(c, fieldErrors)
		return true
	}

	return false
}
