package postgres

import (
	"auth/internal/core/domain"
	"auth/internal/core/ports/output"
	"context"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userConsentRepository struct {
	pool *pgxpool.Pool
}

func NewUserConsentRepository(pool *pgxpool.Pool) output.UserConsentRepository {
	return &userConsentRepository{pool: pool}
}

func (r *userConsentRepository) Save(ctx context.Context, userConsent *domain.UserConsent) (*domain.UserConsent, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("save_user_consent"),
		logger.DBTable("user_consent"),
	)

	log.Debug("inserting user consent")

	query := `
        INSERT INTO user_consents (id, user_id, type, accepted, version, ip_address, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, user_id, type, accepted, version, ip_address, created_at
    `

	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx, query,
		userConsent.ID,
		userConsent.UserID,
		userConsent.Type,
		userConsent.Accepted,
		userConsent.Version,
		userConsent.IPAddress,
		userConsent.CreatedAt,
	)

	saved, err := scanUserConsent(row)
	if err != nil {
		log.Error("failed to insert user consent", logger.Err(err))
		return nil, err
	}

	log.Info("user consent inserted", logger.UserID(saved.UserID.String()))
	return saved, nil
}

func scanUserConsent(row pgx.Row) (*domain.UserConsent, error) {
	var uc domain.UserConsent
	err := row.Scan(
		&uc.ID,
		&uc.UserID,
		&uc.Type,
		&uc.Accepted,
		&uc.Version,
		&uc.IPAddress,
		&uc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &uc, nil
}
