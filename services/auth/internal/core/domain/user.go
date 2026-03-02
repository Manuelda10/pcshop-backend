package domain

import (
	"time"

	"github.com/google/uuid"
)

// Provider representa el origen del registro del usuario.
type Provider string

const (
	ProviderLocal  Provider = "local"
	ProviderGoogle Provider = "google"
)

// User es la entidad central del dominio de autenticación.
// No conoce nada de base de datos, HTTP ni ningún framework.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string // vacío si el proveedor es Google
	FullName     string
	Provider     Provider
	ProviderID   string // Google sub ID, vacío si es local
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewLocalUser crea un usuario de registro manual.
// Recibe el hash ya procesado — el hashing es responsabilidad del service.
func NewLocalUser(email, passwordHash, fullName string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Provider:     ProviderLocal,
		IsVerified:   false, // requiere verificación de email
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewGoogleUser crea un usuario proveniente de OAuth2 con Google.
// Se marca como verificado porque Google ya validó el email.
func NewGoogleUser(email, fullName, providerID string) *User {
	now := time.Now()
	return &User{
		ID:         uuid.New(),
		Email:      email,
		FullName:   fullName,
		Provider:   ProviderGoogle,
		ProviderID: providerID,
		IsVerified: true, // Google garantiza que el email es válido
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// IsLocal retorna true si el usuario se registró manualmente.
func (u *User) IsLocal() bool {
	return u.Provider == ProviderLocal
}

// CanLogin verifica si el usuario puede iniciar sesión.
// Un usuario local no verificado no puede loguearse.
func (u *User) CanLogin() error {
	if u.IsLocal() && !u.IsVerified {
		return ErrEmailNotVerified
	}
	return nil
}

// Verify marca el usuario como verificado.
func (u *User) Verify() {
	u.IsVerified = true
	u.UpdatedAt = time.Now()
}
