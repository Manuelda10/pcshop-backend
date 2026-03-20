package input

import (
	"auth/internal/core/domain"
	"context"
)

// AuthService define los casos de uso del servicio de autenticación.
// Es el "puerto de entrada" — el handler solo conoce esta interfaz, nunca la implementación.
//
// Implementado en: internal/core/services/auth_service.go
type AuthService interface {
	// Register crea un nuevo usuario local y envía email de verificación.
	Register(ctx context.Context, cmd RegisterCommand) (*domain.User, error)

	/*
		// VerifyEmail confirma el código enviado al email del usuario.
		VerifyEmail(ctx context.Context, cmd VerifyEmailCommand) error

		// Login autentica un usuario local y retorna los tokens JWT.
		Login(ctx context.Context, cmd LoginCommand) (*AuthTokens, error)

		// Logout invalida el refresh token del usuario.
		Logout(ctx context.Context, cmd LogoutCommand) error

		// RefreshToken valida el refresh token actual, lo invalida y emite uno nuevo.
		RefreshToken(ctx context.Context, cmd RefreshTokenCommand) (*AuthTokens, error)

		// GoogleCallback procesa el callback de Google OAuth2 y retorna tokens.
		// Si el usuario no existe lo crea automáticamente.
		GoogleCallback(ctx context.Context, cmd GoogleCallbackCommand) (*AuthTokens, error)

		// Me retorna la información del usuario autenticado.
		Me(ctx context.Context, userID string) (*domain.User, error)*/
}

// -------------------------
// Commands (input de cada caso de uso)
// -------------------------

// TODO: modificar servicio para aceptar más parámetros. Además, validarlos y hacerlos obligatorios. Además, cambiar el registro para usar @FirstName y @LastName
type RegisterCommand struct {
	Email          string
	Password       string
	FirstName      string
	LastName       string
	DocumentType   string
	DocumentNumber string
	PhoneNumber    string
	Consents       ConsentsCommand
	IPAddress      string
}

type ConsentsCommand struct {
	PrivacyPolicy bool
	DataCampaign  bool
}

type VerifyEmailCommand struct {
	UserID string
	Code   string
}

type LoginCommand struct {
	Email    string
	Password string
}

type LogoutCommand struct {
	RefreshToken string
}

type RefreshTokenCommand struct {
	RefreshToken string
}

type GoogleCallbackCommand struct {
	Code  string // authorization code que retorna Google
	State string // para validar CSRF
}

// -------------------------
// Responses
// -------------------------

// AuthTokens agrupa los tokens que se retornan al cliente tras login/registro.
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int // segundos hasta que expira el access token
}
