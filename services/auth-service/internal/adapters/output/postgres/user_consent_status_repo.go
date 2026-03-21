package postgres

import (
	"auth-service/internal/core/domain"
	"auth-service/internal/core/ports/output"
	"context"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userConsentStatusRepository struct {
	pool *pgxpool.Pool
}

func NewUserConsentStatusRepository(pool *pgxpool.Pool) output.UserConsentStatusRepository {
	return &userConsentStatusRepository{pool: pool}
}

func (r *userConsentStatusRepository) Save(ctx context.Context, userConsentStatus *domain.UserConsentStatus) (*domain.UserConsentStatus, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("save_user_consent_status"),
		logger.DBTable("user_consent_status"),
	)

	log.Debug("inserting user consent status")

	query := `
        INSERT INTO user_consent_status (user_id, consented_privacy_policy, consented_data_campaign, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING user_id, consented_privacy_policy, consented_data_campaign, created_at, updated_at
    `

	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx, query,
		userConsentStatus.UserID,
		userConsentStatus.ConsentedPrivacyPolicy,
		userConsentStatus.ConsentedDataCampaign,
		userConsentStatus.CreatedAt,
		userConsentStatus.UpdatedAt,
	)

	saved, err := scanUserConsentStatus(row)
	if err != nil {
		log.Error("failed to insert user consent status", logger.Err(err))
		return nil, err
	}

	log.Info("user consent status inserted", logger.UserID(saved.UserID.String()))
	return saved, nil
}

func scanUserConsentStatus(row pgx.Row) (*domain.UserConsentStatus, error) {
	var ucs domain.UserConsentStatus
	err := row.Scan(
		&ucs.UserID,
		&ucs.ConsentedPrivacyPolicy,
		&ucs.ConsentedDataCampaign,
		&ucs.CreatedAt,
		&ucs.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &ucs, nil
}
