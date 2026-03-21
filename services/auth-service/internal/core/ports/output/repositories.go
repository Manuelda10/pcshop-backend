package output

import (
	"auth-service/internal/core/domain"
	"context"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	/*
		// FindByID busca un usuario por su ID.
		FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

		// FindByEmail busca un usuario por email.


		// FindByProviderID busca un usuario por su Google sub ID.
		FindByProviderID(ctx context.Context, provider domain.Provider, providerID string) (*domain.User, error)

		// Update actualiza los campos modificables de un usuario.
		Update(ctx context.Context, user *domain.User) (*domain.User, error)*/
}

type UserConsentStatusRepository interface {
	Save(ctx context.Context, userConsentStatus *domain.UserConsentStatus) (*domain.UserConsentStatus, error)
}

type UserConsentRepository interface {
	Save(ctx context.Context, userConsent *domain.UserConsent) (*domain.UserConsent, error)
}

type TokenRepository interface {
	// SaveRefreshToken persiste un nuevo refresh token.
	SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error

	/*
		// FindRefreshToken busca un refresh token por su hash.
		FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)

		// DeleteRefreshToken elimina un refresh token (logout).
		DeleteRefreshToken(ctx context.Context, tokenHash string) error

		// DeleteAllUserRefreshTokens elimina todos los tokens de un usuario.
		// Útil para "cerrar todas las sesiones".
		DeleteAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
	*/
	// SaveEmailVerification persiste un código de verificación de email.
	SaveEmailVerification(ctx context.Context, ev *domain.EmailVerification) error

	/*
		// FindEmailVerification busca un código de verificación por userID.
		FindEmailVerification(ctx context.Context, userID uuid.UUID) (*domain.EmailVerification, error)

		// DeleteEmailVerification elimina el código tras verificación exitosa.
		DeleteEmailVerification(ctx context.Context, userID uuid.UUID) error
	*/
}

type EmailSender interface {
	// SendVerificationEmail envía el código de verificación al usuario recién registrado.
	SendVerificationEmail(ctx context.Context, toEmail, toName, code string) error

	// SendPasswordResetEmail envía el link/código de reseteo de contraseña.
	// (reservado para implementación futura)
	SendPasswordResetEmail(ctx context.Context, toEmail, toName, code string) error
}
