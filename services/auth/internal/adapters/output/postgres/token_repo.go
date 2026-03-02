package postgres

import (
	"auth-service/internal/core/domain"
	"auth-service/internal/core/ports/output"
	"context"
	"errors"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tokenRepository struct {
	db *pgxpool.Pool
}

// NewTokenRepository crea una nueva instancia del repositorio de tokens.
func NewTokenRepository(db *pgxpool.Pool) output.TokenRepository {
	return &tokenRepository{db: db}
}

// -------------------------
// Refresh Tokens
// -------------------------

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("save_refresh_token"),
		logger.DBTable("refresh_tokens"),
		logger.UserID(token.UserID.String()),
	)

	log.Debug("inserting refresh token")

	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
	)
	if err != nil {
		log.Error("failed to insert refresh token", logger.Err(err))
		return err
	}

	log.Debug("refresh token saved")
	return nil
}

func (r *tokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("find_refresh_token"),
		logger.DBTable("refresh_tokens"),
	)

	log.Debug("querying refresh token by hash")

	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens WHERE token_hash = $1
	`

	var rt domain.RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("refresh token not found")
			return nil, domain.ErrTokenNotFound
		}
		log.Error("failed to query refresh token", logger.Err(err))
		return nil, err
	}

	log.Debug("refresh token found", logger.UserID(rt.UserID.String()))
	return &rt, nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("delete_refresh_token"),
		logger.DBTable("refresh_tokens"),
	)

	log.Debug("deleting refresh token")

	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`

	result, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		log.Error("failed to delete refresh token", logger.Err(err))
		return err
	}

	if result.RowsAffected() == 0 {
		log.Warn("no refresh token found to delete")
		return domain.ErrTokenNotFound
	}

	log.Debug("refresh token deleted")
	return nil
}

func (r *tokenRepository) DeleteAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("delete_all_user_tokens"),
		logger.DBTable("refresh_tokens"),
		logger.UserID(userID.String()),
	)

	log.Debug("deleting all refresh tokens for user")

	query := `DELETE FROM refresh_tokens WHERE user_id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		log.Error("failed to delete user tokens", logger.Err(err))
		return err
	}

	log.Info("all user tokens deleted", logger.Output(map[string]int64{"rows_affected": result.RowsAffected()}))
	return nil
}

// -------------------------
// Email Verifications
// -------------------------

func (r *tokenRepository) SaveEmailVerification(ctx context.Context, ev *domain.EmailVerification) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("save_email_verification"),
		logger.DBTable("email_verifications"),
		logger.UserID(ev.UserID.String()),
	)

	log.Debug("inserting email verification code")

	// Upsert — si ya existe uno para este usuario, lo reemplaza
	query := `
		INSERT INTO email_verifications (id, user_id, code_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE
		SET id = EXCLUDED.id, code_hash = EXCLUDED.code_hash,
		    expires_at = EXCLUDED.expires_at, created_at = EXCLUDED.created_at
	`

	_, err := r.db.Exec(ctx, query,
		ev.ID,
		ev.UserID,
		ev.CodeHash,
		ev.ExpiresAt,
		ev.CreatedAt,
	)
	if err != nil {
		log.Error("failed to save email verification", logger.Err(err))
		return err
	}

	log.Debug("email verification saved")
	return nil
}

func (r *tokenRepository) FindEmailVerification(ctx context.Context, userID uuid.UUID) (*domain.EmailVerification, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("find_email_verification"),
		logger.DBTable("email_verifications"),
		logger.UserID(userID.String()),
	)

	log.Debug("querying email verification")

	query := `
		SELECT id, user_id, code_hash, expires_at, created_at
		FROM email_verifications WHERE user_id = $1
	`

	var ev domain.EmailVerification
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&ev.ID,
		&ev.UserID,
		&ev.CodeHash,
		&ev.ExpiresAt,
		&ev.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("email verification not found")
			return nil, domain.ErrVerificationCodeNotFound
		}
		log.Error("failed to query email verification", logger.Err(err))
		return nil, err
	}

	log.Debug("email verification found")
	return &ev, nil
}

func (r *tokenRepository) DeleteEmailVerification(ctx context.Context, userID uuid.UUID) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("delete_email_verification"),
		logger.DBTable("email_verifications"),
		logger.UserID(userID.String()),
	)

	log.Debug("deleting email verification")

	query := `DELETE FROM email_verifications WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		log.Error("failed to delete email verification", logger.Err(err))
		return err
	}

	log.Debug("email verification deleted")
	return nil
}
