package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func OK(c *fiber.Ctx, code ResponseCode, message string, data any) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Code:    string(code),
		Message: message,
		Data:    data,
		Errors:  nil,
		Meta:    nil,
	})
}

func OKWithMeta(c *fiber.Ctx, code ResponseCode, message string, data, meta any) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Code:    string(code),
		Message: message,
		Data:    data,
		Errors:  nil,
		Meta:    meta,
	})
}

func Created(c *fiber.Ctx, code ResponseCode, message string, data any) error {
	return c.Status(http.StatusCreated).JSON(Response{
		Success: true,
		Code:    string(code),
		Message: message,
		Data:    data,
		Errors:  nil,
		Meta:    nil,
	})
}

func NoContent(c *fiber.Ctx, code ResponseCode, message string) error {
	return c.Status(http.StatusNoContent).JSON(Response{
		Success: true,
		Code:    string(code),
		Message: message,
		Data:    nil,
		Errors:  nil,
		Meta:    nil,
	})
}

func BadRequest(c *fiber.Ctx, code ResponseCode, message string) error {
	return c.Status(http.StatusBadRequest).JSON(Response{
		Success: false,
		Code:    string(code),
		Message: message,
		Data:    nil,
		Errors:  nil,
		Meta:    nil,
	})
}

func ValidationError(c *fiber.Ctx, errors []FieldError) error {
	return c.Status(http.StatusUnprocessableEntity).JSON(Response{
		Success: false,
		Code:    string(CodeValidationError),
		Message: "Error de validación en los datos enviados",
		Data:    nil,
		Errors:  errors,
		Meta:    nil,
	})
}

func Unauthorized(c *fiber.Ctx, code ResponseCode, message string) error {
	return c.Status(http.StatusUnauthorized).JSON(Response{
		Success: false,
		Code:    string(code),
		Message: message,
		Data:    nil,
		Errors:  nil,
		Meta:    nil,
	})
}

func Conflict(c *fiber.Ctx, code ResponseCode, message string) error {
	return c.Status(http.StatusConflict).JSON(Response{
		Success: false,
		Code:    string(code),
		Message: message,
		Data:    nil,
		Errors:  nil,
		Meta:    nil,
	})
}

func NotFound(c *fiber.Ctx, code ResponseCode, message string) error {
	return c.Status(http.StatusNotFound).JSON(Response{
		Success: false,
		Code:    string(code),
		Message: message,
		Data:    nil,
		Errors:  nil,
		Meta:    nil,
	})
}

func InternalError(c *fiber.Ctx) error {
	return c.Status(http.StatusInternalServerError).JSON(Response{
		Success: false,
		Code:    string(CodeInternalError),
		Message: "Ocurrió un error inesperado, intente más tarde",
		Data:    nil,
		Errors:  nil,
		Meta:    nil,
	})
}
