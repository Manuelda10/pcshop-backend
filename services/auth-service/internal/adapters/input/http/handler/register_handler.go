package handler

import (
	"auth-service/internal/adapters/input/http/dto"
	"auth-service/internal/core/ports/input"

	sharedhttp "github.com/Manuelda10/pcshop-backend/shared/http"
	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/Manuelda10/pcshop-backend/shared/validation"
	"github.com/gofiber/fiber/v2"
)

type RegisterHandler struct {
	Base
	service input.AuthService
}

func NewRegisterHandler(service input.AuthService, v *validation.Validator) *RegisterHandler {
	return &RegisterHandler{
		Base:    Base{Validator: v},
		service: service,
	}
}

func (h *RegisterHandler) Handle(c *fiber.Ctx) error {
	log := logger.FromCtx(c).With(logger.Operation("register"))

	var req dto.RegisterRequest
	if failed := parseAndValidate(c, &h.Base, &req); failed {
		return nil
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
		return handleServiceError(c, err)
	}

	log.Info("register handler completed", logger.UserID(user.ID.String()))
	return sharedhttp.Created(c, CodeUserCreated, "Usuario creado exitosamente", dto.ToUserResponse(user))
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
