package domain

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderLocal  Provider = "local"
	ProviderGoogle Provider = "google"
)

type User struct {
	ID             uuid.UUID
	Email          string
	PasswordHash   string
	FirstName      string
	LastName       string
	DocumentType   string
	DocumentNumber string
	PhoneNumber    string
	Provider       Provider
	ProviderID     *string
	IsVerified     bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NewUserParams struct {
	Email          string
	PasswordHash   string
	FirstName      string
	LastName       string
	DocumentType   string
	DocumentNumber string
	PhoneNumber    string
	ProviderID     string
}

func NewLocalUser(u NewUserParams) *User {
	now := time.Now()
	return &User{
		ID:             uuid.New(),
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		DocumentType:   u.DocumentType,
		DocumentNumber: u.DocumentNumber,
		PhoneNumber:    u.PhoneNumber,
		Provider:       ProviderLocal,
		IsVerified:     false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func NewGoogleUser(u NewUserParams) *User {
	now := time.Now()
	return &User{
		ID:             uuid.New(),
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		DocumentType:   u.DocumentType,
		DocumentNumber: u.DocumentNumber,
		PhoneNumber:    u.PhoneNumber,
		Provider:       ProviderGoogle,
		ProviderID:     &u.ProviderID,
		IsVerified:     true,
		CreatedAt:      now,
		UpdatedAt:      now,
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
