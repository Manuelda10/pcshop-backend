package domain

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken representa el token de refresco almacenado en base de datos.
// Nunca guardamos el token plano, solo su hash.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string // bcrypt hash del token real
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewRefreshToken crea un nuevo refresh token para un usuario.
func NewRefreshToken(userID uuid.UUID, tokenHash string, expiry time.Duration) *RefreshToken {
	now := time.Now()
	return &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(expiry),
		CreatedAt: now,
	}
}

// IsExpired verifica si el token ha expirado.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// -------------------------

// EmailVerification representa el código de verificación enviado por email.
type EmailVerification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CodeHash  string // hash del código de 6 dígitos enviado al usuario
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewEmailVerification crea un nuevo registro de verificación de email.
func NewEmailVerification(userID uuid.UUID, codeHash string) *EmailVerification {
	now := time.Now()
	return &EmailVerification{
		ID:        uuid.New(),
		UserID:    userID,
		CodeHash:  codeHash,
		ExpiresAt: now.Add(24 * time.Hour), // el código expira en 24 horas
		CreatedAt: now,
	}
}

// IsExpired verifica si el código de verificación ha expirado.
func (ev *EmailVerification) IsExpired() bool {
	return time.Now().After(ev.ExpiresAt)
}
