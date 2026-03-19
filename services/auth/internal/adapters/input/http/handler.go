package http

import (
	"auth-service/config"
	"auth-service/internal/adapters/input/http/dto"
	"auth-service/internal/core/domain"
	"auth-service/internal/core/ports/input"
	"errors"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// AuthHandler contiene los handlers HTTP para el servicio de autenticación.
// Solo conoce el port de entrada (input.AuthService) — nunca la implementación.
type AuthHandler struct {
	service      input.AuthService
	googleConfig *oauth2.Config
}

// NewAuthHandler construye el handler con sus dependencias.
func NewAuthHandler(service input.AuthService, cfg config.Config) *AuthHandler {
	googleCfg := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &AuthHandler{
		service:      service,
		googleConfig: googleCfg,
	}
}

// -------------------------
// Register
// -------------------------

// Register godoc
// @Summary      Registro de nuevo usuario
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RegisterRequest true "Datos de registro"
// @Success      201  {object} dto.UserResponse
// @Failure      400  {object} dto.ErrorResponse
// @Failure      409  {object} dto.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("register"))

	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		log.Warn("invalid request body", logger.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	log.Debug("handler received register request", logger.Input(map[string]string{
		"email":          req.Email,
		"documentNumber": req.DocumentNumber,
		"ip_address":     c.IP(),
	}))

	user, err := h.service.Register(c.UserContext(), input.RegisterCommand{
		Email:          req.Email,
		Password:       req.Password,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		PhoneNumber:    req.PhoneNumber,
		Consents: input.ConsentsCommand{
			PrivacyPolicy: req.Consents.PrivacyPolicy,
			DataCampaign:  req.Consents.DataCampaign,
		},
		IPAddress: c.IP(),
	})
	if err != nil {
		return h.handleServiceError(c, err)
	}

	log.Info("register handler completed", logger.UserID(user.ID.String()))
	return c.Status(fiber.StatusCreated).JSON(dto.ToUserResponse(user))
}

/*
// -------------------------
// VerifyEmail
// -------------------------

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("verify_email"))

	var req dto.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		log.Warn("invalid request body", logger.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	log.Debug("handler received verify email request", logger.UserID(req.UserID))

	if err := h.service.VerifyEmail(c.UserContext(), input.VerifyEmailCommand{
		UserID: req.UserID,
		Code:   req.Code,
	}); err != nil {
		return h.handleServiceError(c, err)
	}

	log.Info("email verified", logger.UserID(req.UserID))
	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "email verified successfully — you can now log in",
	})
}

// -------------------------
// Login
// -------------------------

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("login"))

	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Warn("invalid request body", logger.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	log.Debug("handler received login request", logger.Email(req.Email))

	tokens, err := h.service.Login(c.UserContext(), input.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Buscar el usuario para incluirlo en la respuesta
	user, err := h.service.Me(c.UserContext(), extractUserIDFromToken(tokens.AccessToken))
	if err != nil {
		// No bloqueante — si falla igual retornamos los tokens
		log.Warn("could not fetch user after login", logger.Err(err))
		return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TokenType:    tokens.TokenType,
			ExpiresIn:    tokens.ExpiresIn,
		})
	}

	log.Info("login handler completed", logger.UserID(user.ID.String()))
	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
		User:         dto.ToUserResponse(user),
	})
}

// -------------------------
// Logout
// -------------------------

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("logout"))

	var req dto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		log.Warn("invalid request body", logger.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	if err := h.service.Logout(c.UserContext(), input.LogoutCommand{
		RefreshToken: req.RefreshToken,
	}); err != nil {
		return h.handleServiceError(c, err)
	}

	log.Info("logout handler completed")
	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "logged out successfully",
	})
}

// -------------------------
// RefreshToken
// -------------------------

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("refresh_token"))

	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		log.Warn("invalid request body", logger.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	tokens, err := h.service.RefreshToken(c.UserContext(), input.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return h.handleServiceError(c, err)
	}

	log.Info("token refreshed")
	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// -------------------------
// Google OAuth2
// -------------------------

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	// Generar state para prevenir CSRF
	state := generateOAuthState()
	url := h.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Guardar state en cookie temporal para validar en callback
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   300, // 5 minutos
	})

	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("google_callback"))

	// Validar state CSRF
	cookieState := c.Cookies("oauth_state")
	queryState := c.Query("state")
	if cookieState == "" || cookieState != queryState {
		log.Warn("invalid oauth state — possible CSRF")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid state parameter",
			Code:  "INVALID_OAUTH_STATE",
		})
	}

	code := c.Query("code")
	if code == "" {
		log.Warn("missing authorization code in google callback")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "missing authorization code",
			Code:  "MISSING_AUTH_CODE",
		})
	}

	log.Debug("processing google oauth code exchange")

	// Intercambiar code por token de Google
	googleToken, err := h.googleConfig.Exchange(c.UserContext(), code)
	if err != nil {
		log.Error("failed to exchange google code", logger.Err(err))
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "google authentication failed",
			Code:  "GOOGLE_AUTH_FAILED",
		})
	}

	// Obtener perfil del usuario desde Google
	googleUser, err := getGoogleUserInfo(c.UserContext(), googleToken)
	if err != nil {
		log.Error("failed to fetch google user info", logger.Err(err))
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "failed to fetch google profile",
			Code:  "GOOGLE_PROFILE_FAILED",
		})
	}

	log.Debug("google user fetched", logger.Email(googleUser.Email))

	tokens, err := h.service.GoogleCallback(c.UserContext(), input.GoogleCallbackCommand{
		Code:  code,
		State: queryState,
	})
	if err != nil {
		return h.handleServiceError(c, err)
	}

	_ = googleUser // usado en la implementación completa del service

	log.Info("google login completed")
	return c.Status(fiber.StatusOK).JSON(tokens)
}

// -------------------------
// Me
// -------------------------

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("me"))

	// El userID lo inyecta el middleware JWT (ver routes.go)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		log.Warn("userID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
			Code:  "UNAUTHORIZED",
		})
	}

	user, err := h.service.Me(c.UserContext(), userID)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	log.Debug("me handler completed", logger.UserID(userID))
	return c.Status(fiber.StatusOK).JSON(dto.ToUserResponse(user))
}

*/
// -------------------------
// Error mapper
// -------------------------

// handleServiceError mapea los errores de dominio a respuestas HTTP apropiadas.
// Centralizado aquí para no repetir lógica en cada handler.
func (h *AuthHandler) handleServiceError(c *fiber.Ctx, err error) error {
	log := logger.FromCtx(c)

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: err.Error(), Code: "USER_NOT_FOUND"})

	case errors.Is(err, domain.ErrUserAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: err.Error(), Code: "USER_ALREADY_EXISTS"})

	case errors.Is(err, domain.ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: err.Error(), Code: "INVALID_CREDENTIALS"})

	case errors.Is(err, domain.ErrEmailNotVerified):
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{Error: err.Error(), Code: "EMAIL_NOT_VERIFIED"})

	case errors.Is(err, domain.ErrInvalidToken),
		errors.Is(err, domain.ErrTokenNotFound):
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: err.Error(), Code: "INVALID_TOKEN"})

	case errors.Is(err, domain.ErrExpiredToken):
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: err.Error(), Code: "TOKEN_EXPIRED"})

	case errors.Is(err, domain.ErrInvalidVerificationCode),
		errors.Is(err, domain.ErrExpiredVerificationCode),
		errors.Is(err, domain.ErrVerificationCodeNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error(), Code: "INVALID_VERIFICATION_CODE"})

	case errors.Is(err, domain.ErrProviderMismatch):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: err.Error(), Code: "PROVIDER_MISMATCH"})

	default:
		log.Error("unhandled service error", logger.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}
}
