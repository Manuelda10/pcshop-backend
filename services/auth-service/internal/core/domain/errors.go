package domain

import "errors"

// Errores de dominio — representan situaciones de negocio, no errores técnicos.
// El service los retorna y el handler los mapea a códigos HTTP.
//
// Uso:
//   if errors.Is(err, domain.ErrUserNotFound) {
//       return c.Status(404).JSON(...)
//   }

var (
	// Usuario
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already registered")
	ErrDocumentAlreadyExists = errors.New("document already registered")
	ErrUserAlreadyExists     = errors.New("user already exists with this email")
	ErrEmailNotVerified      = errors.New("email not verified — check your inbox")
	ErrInvalidCredentials    = errors.New("invalid email or password")

	// Tokens
	ErrInvalidToken  = errors.New("invalid or malformed token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrTokenNotFound = errors.New("token not found")

	// Verificación de email
	ErrInvalidVerificationCode  = errors.New("invalid verification code")
	ErrExpiredVerificationCode  = errors.New("verification code has expired")
	ErrVerificationCodeNotFound = errors.New("verification code not found")

	// Google OAuth
	ErrGoogleAuthFailed = errors.New("google authentication failed")
	ErrProviderMismatch = errors.New("email already registered with a different provider")

	// Genérico
	ErrInternalServer = errors.New("internal server error")
)
