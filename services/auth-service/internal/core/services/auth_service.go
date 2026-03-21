package services

import (
	"auth-service/config"
	"auth-service/internal/core/domain"
	"auth-service/internal/core/ports/input"
	"auth-service/internal/core/ports/output"
	"context"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// authService implementa input.AuthService.
// Solo conoce los ports de output — nunca los adaptadores concretos.
type authService struct {
	txManager         output.TxManager
	userRepo          output.UserRepository
	consentStatusRepo output.UserConsentStatusRepository
	consentRepo       output.UserConsentRepository
	tokenRepo         output.TokenRepository
	//emailSender       output.EmailSender
	cfg config.Config
}

func NewAuthService(
	txManager output.TxManager,
	userRepo output.UserRepository,
	consentStatusRepo output.UserConsentStatusRepository,
	consentRepo output.UserConsentRepository,
	tokenRepo output.TokenRepository,
	//emailSender output.EmailSender,
	cfg config.Config,
) input.AuthService {
	return &authService{
		txManager:         txManager,
		userRepo:          userRepo,
		consentStatusRepo: consentStatusRepo,
		consentRepo:       consentRepo,
		tokenRepo:         tokenRepo,
		//emailSender:       emailSender,
		cfg: cfg,
	}
}

// -------------------------
// Register
// -------------------------

func (s *authService) Register(ctx context.Context, cmd input.RegisterCommand) (*domain.User, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("register"),
		logger.Email(cmd.Email),
	)

	log.Debug("registering new user")

	existing, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		log.Error("failed to check existing email", logger.Err(err))
		return nil, domain.ErrInternalServer
	}
	if existing != nil {
		log.Warn("email already registered")
		return nil, domain.ErrEmailAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", logger.Err(err))
		return nil, domain.ErrInternalServer
	}

	// Crear entidad de dominio
	user := domain.NewLocalUser(domain.NewUserParams{
		Email:          cmd.Email,
		PasswordHash:   string(passwordHash),
		FirstName:      cmd.FirstName,
		LastName:       cmd.LastName,
		DocumentType:   cmd.DocumentType,
		DocumentNumber: cmd.DocumentNumber,
		PhoneNumber:    cmd.PhoneNumber,
	})
	userConsentStatus := domain.NewUserConsentStatus(user.ID, cmd.Consents.PrivacyPolicy, cmd.Consents.DataCampaign)
	consentHistory := []*domain.UserConsent{
		domain.NewUserConsent(domain.NewUserConsentParams{
			UserID:    user.ID,
			Type:      domain.ConsentPrivacyPolicy,
			Accepted:  cmd.Consents.PrivacyPolicy,
			Version:   s.cfg.Legal.ConsentPrivacyVersion,
			IPAddress: cmd.IPAddress,
		}),
		domain.NewUserConsent(domain.NewUserConsentParams{
			UserID:    user.ID,
			Type:      domain.ConsentDataCampaign,
			Accepted:  cmd.Consents.DataCampaign,
			Version:   s.cfg.Legal.ConsentDataCampaignVersion,
			IPAddress: cmd.IPAddress,
		}),
	}

	var saved *domain.User

	err = s.txManager.RunInTx(ctx, func(ctx context.Context) error {
		var err error

		saved, err = s.userRepo.Save(ctx, user)
		if err != nil {
			log.Error("failed to save user", logger.Err(err))
			return err
		}

		_, err = s.consentStatusRepo.Save(ctx, userConsentStatus)
		if err != nil {
			log.Error("failed to save consent status", logger.Err(err))
			return err
		}

		for _, consent := range consentHistory {
			_, err = s.consentRepo.Save(ctx, consent)
			if err != nil {
				log.Error("failed to save consent history", logger.Err(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	// Generar código de verificación (6 dígitos)
	code, codeHash, err := generateVerificationCode()
	if err != nil {
		log.Error("failed to generate verification code", logger.Err(err))
		return nil, domain.ErrInternalServer
	}
	//TODO: Eliminar cuando se implemente el envío de correo
	log.Info("otp code sent", logger.Input(code))

	ev := domain.NewEmailVerification(saved.ID, codeHash)
	if err := s.tokenRepo.SaveEmailVerification(ctx, ev); err != nil {
		log.Error("failed to save email verification", logger.Err(err))
		return nil, domain.ErrInternalServer
	}

	// Enviar email (no bloqueante — si falla loguea pero no rompe el registro) Comentado hasta implementación correcta.
	/*
		if err := s.emailSender.SendVerificationEmail(ctx, saved.Email, saved.FullName, code); err != nil {
			log.Warn("failed to send verification email — user registered but email not sent", logger.Err(err))
		}*/

	log.Info("user registered successfully", logger.UserID(saved.ID.String()))
	return saved, nil
}

// -------------------------
// VerifyEmail
// -------------------------
/*
func (s *authService) VerifyEmail(ctx context.Context, cmd input.VerifyEmailCommand) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("verify_email"),
		logger.UserID(cmd.UserID),
	)

	log.Debug("verifying email")

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return domain.ErrInvalidToken
	}

	ev, err := s.tokenRepo.FindEmailVerification(ctx, userID)
	if err != nil {
		log.Warn("verification code not found", logger.Err(err))
		return domain.ErrVerificationCodeNotFound
	}

	if ev.IsExpired() {
		log.Warn("verification code expired")
		return domain.ErrExpiredVerificationCode
	}

	// Comparar código hasheado
	if err := bcrypt.CompareHashAndPassword([]byte(ev.CodeHash), []byte(cmd.Code)); err != nil {
		log.Warn("invalid verification code")
		return domain.ErrInvalidVerificationCode
	}

	// Marcar usuario como verificado
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Error("user not found after verification", logger.Err(err))
		return domain.ErrInternalServer
	}

	user.Verify()
	if _, err := s.userRepo.Update(ctx, user); err != nil {
		log.Error("failed to update user verification status", logger.Err(err))
		return domain.ErrInternalServer
	}

	// Limpiar el código usado
	_ = s.tokenRepo.DeleteEmailVerification(ctx, userID)

	log.Info("email verified successfully")
	return nil
}

// -------------------------
// Login
// -------------------------

func (s *authService) Login(ctx context.Context, cmd input.LoginCommand) (*input.AuthTokens, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("login"),
		logger.Email(cmd.Email),
	)

	log.Debug("processing login")

	user, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		// No revelar si el email existe o no — siempre "invalid credentials"
		log.Warn("user not found during login")
		return nil, domain.ErrInvalidCredentials
	}

	if err := user.CanLogin(); err != nil {
		log.Warn("user cannot login", logger.Err(err))
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cmd.Password)); err != nil {
		log.Warn("invalid password")
		return nil, domain.ErrInvalidCredentials
	}

	tokens, err := s.issueTokens(ctx, user)
	if err != nil {
		log.Error("failed to issue tokens", logger.Err(err))
		return nil, domain.ErrInternalServer
	}

	log.Info("login successful", logger.UserID(user.ID.String()))
	return tokens, nil
}

// -------------------------
// Logout
// -------------------------

func (s *authService) Logout(ctx context.Context, cmd input.LogoutCommand) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("logout"),
	)

	log.Debug("processing logout")

	tokenHash, err := hashToken(cmd.RefreshToken)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if err := s.tokenRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		log.Warn("refresh token not found on logout", logger.Err(err))
		// No es un error crítico — el token quizás ya expiró
		return nil
	}

	log.Info("logout successful")
	return nil
}

// -------------------------
// RefreshToken
// -------------------------

func (s *authService) RefreshToken(ctx context.Context, cmd input.RefreshTokenCommand) (*input.AuthTokens, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("refresh_token"),
	)

	log.Debug("rotating refresh token")

	tokenHash, err := hashToken(cmd.RefreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	stored, err := s.tokenRepo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		log.Warn("refresh token not found")
		return nil, domain.ErrTokenNotFound
	}

	if stored.IsExpired() {
		log.Warn("refresh token expired", logger.UserID(stored.UserID.String()))
		_ = s.tokenRepo.DeleteRefreshToken(ctx, tokenHash)
		return nil, domain.ErrExpiredToken
	}

	user, err := s.userRepo.FindByID(ctx, stored.UserID)
	if err != nil {
		log.Error("user not found during token refresh", logger.Err(err))
		return nil, domain.ErrInternalServer
	}

	// Invalidar el token anterior (rotation)
	_ = s.tokenRepo.DeleteRefreshToken(ctx, tokenHash)

	tokens, err := s.issueTokens(ctx, user)
	if err != nil {
		log.Error("failed to issue new tokens", logger.Err(err))
		return nil, domain.ErrInternalServer
	}

	log.Info("token refreshed successfully", logger.UserID(user.ID.String()))
	return tokens, nil
}

// -------------------------
// GoogleCallback
// -------------------------

func (s *authService) GoogleCallback(ctx context.Context, cmd input.GoogleCallbackCommand) (*input.AuthTokens, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("google_callback"),
	)

	log.Debug("processing google oauth callback")

	// La lógica de intercambio de code → userInfo se hace en el handler/adapter
	// porque requiere el oauth2.Config de Google. El service recibe el resultado.
	// Ver: adapters/input/http/handler.go → googleCallbackHandler
	//
	// Este método se llama después de obtener el perfil de Google.
	// Por eso el command lleva el code crudo — el adapter lo resuelve antes de llegar aquí.
	// (implementación completa en el adapter HTTP)

	log.Warn("google callback not fully implemented in service — handled by adapter")
	return nil, nil
}

// -------------------------
// Me
// -------------------------

func (s *authService) Me(ctx context.Context, userID string) (*domain.User, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("service"),
		logger.Operation("me"),
		logger.UserID(userID),
	)

	log.Debug("fetching current user")

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Warn("user not found", logger.Err(err))
		return nil, domain.ErrUserNotFound
	}

	log.Debug("user fetched", logger.Output(map[string]string{"email": user.Email}))
	return user, nil
}*/

// -------------------------
// Helpers privados
// -------------------------

// issueTokens genera un par access + refresh token y persiste el refresh.
func (s *authService) issueTokens(ctx context.Context, user *domain.User) (*input.AuthTokens, error) {
	// Access token (JWT firmado)
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Refresh token (random string, se guarda hasheado)
	refreshToken, refreshHash, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	rt := domain.NewRefreshToken(user.ID, refreshHash, s.cfg.JWT.RefreshTokenExpiry)
	if err := s.tokenRepo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return &input.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.cfg.JWT.AccessTokenExpiry.Seconds()),
	}, nil
}

// JWTClaims define el payload del access token.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *authService) generateAccessToken(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWT.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.AccessSecret))
}

// generateRandomToken genera un token aleatorio seguro (32 bytes hex)
// y su hash SHA-256 para almacenar en DB.
// SHA-256 aquí porque necesitamos buscar por hash — bcrypt no sirve para eso
// ya que genera un salt distinto cada vez (no es determinístico).
func generateRandomToken() (plain, hash string, err error) {
	b := make([]byte, 32)
	if _, err = crand.Read(b); err != nil {
		return
	}
	plain = fmt.Sprintf("%x", b)
	hash = hashSHA256(plain)
	return
}

// generateVerificationCode genera un código numérico de 6 dígitos
// usando crypto/rand para garantizar aleatoriedad criptográfica,
// y su hash bcrypt para almacenar en DB.
// Usamos bcrypt aquí (no SHA-256) porque el código es corto y predecible —
// bcrypt con su salt lo protege de ataques de fuerza bruta sobre la DB.
func generateVerificationCode() (code, hash string, err error) {
	// Generar número entre 0 y 999999 de forma criptográficamente segura
	max := big.NewInt(1_000_000)
	n, randErr := crand.Int(crand.Reader, max)
	if randErr != nil {
		err = randErr
		return
	}
	code = fmt.Sprintf("%06d", n.Int64())

	hashed, bcryptErr := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if bcryptErr != nil {
		err = bcryptErr
		return
	}
	hash = string(hashed)
	return
}

// hashToken genera el hash SHA-256 de un refresh token para búsqueda en DB.
// SHA-256 es determinístico — el mismo token siempre produce el mismo hash,
// lo que permite buscar: WHERE token_hash = sha256(token_recibido).
func hashToken(token string) (string, error) {
	if token == "" {
		return "", domain.ErrInvalidToken
	}
	return hashSHA256(token), nil
}

// hashSHA256 es el helper base para hashing determinístico.
func hashSHA256(input string) string {
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}
