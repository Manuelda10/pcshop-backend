package postgres

import (
	"auth/internal/core/domain"
	"auth/internal/core/ports/output"
	"context"
	"errors"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) output.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("save_user"),
		logger.DBTable("users"),
	)

	log.Debug("inserting user")

	query := `
        INSERT INTO users (id, email, password_hash, first_name, last_name, document_type, document_number, 
                           phone_number, provider, provider_id, is_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id, email, password_hash, first_name, last_name, document_type, document_number, 
            phone_number, provider, provider_id, is_verified, created_at, updated_at
    `

	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.DocumentType,
		user.DocumentNumber,
		user.PhoneNumber,
		user.Provider,
		user.ProviderID,
		user.IsVerified,
		user.CreatedAt,
		user.UpdatedAt,
	)

	saved, err := scanUser(row)
	if err != nil {
		log.Error("failed to insert user", logger.Err(err))
		return nil, err
	}

	log.Info("user inserted", logger.UserID(saved.ID.String()))
	return saved, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	log := logger.FromContext(ctx).With(
		logger.Layer("repository"),
		logger.Operation("find_user_by_email"),
		logger.DBTable("users"),
		logger.Email(email),
	)

	log.Debug("querying user by email")

	query := `
		SELECT id, email, password_hash, first_name, last_name, document_type, document_number,
			   phone_number, provider, provider_id, is_verified, created_at, updated_at
		FROM users WHERE email = $1
	`

	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx, query, email)
	user, err := scanUser(row)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("user not found by email")
			return nil, domain.ErrUserNotFound
		}
		log.Error("failed to query user by email", logger.Err(err))
		return nil, err
	}

	log.Debug("user found by email")
	return user, nil
}

/*
	func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
		log := logger.FromContext(ctx).With(
			logger.Layer("repository"),
			logger.Operation("find_user_by_id"),
			logger.DBTable("users"),
			logger.UserID(id.String()),
		)

		log.Debug("querying user by id")

		query := `
			SELECT id, email, password_hash, full_name, provider, provider_id, is_verified, created_at, updated_at
			FROM users WHERE id = $1
		`

		row := r.db.QueryRow(ctx, query, id)
		user, err := scanUser(row)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Warn("user not found")
				return nil, domain.ErrUserNotFound
			}
			log.Error("failed to query user", logger.Err(err))
			return nil, err
		}

		log.Debug("user found")
		return user, nil
	}

	func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
		log := logger.FromContext(ctx).With(
			logger.Layer("repository"),
			logger.Operation("find_user_by_email"),
			logger.DBTable("users"),
			logger.Email(email),
		)

		log.Debug("querying user by email")

		query := `
			SELECT id, email, password_hash, full_name, provider, provider_id, is_verified, created_at, updated_at
			FROM users WHERE email = $1
		`

		row := r.db.QueryRow(ctx, query, email)
		user, err := scanUser(row)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Warn("user not found by email")
				return nil, domain.ErrUserNotFound
			}
			log.Error("failed to query user by email", logger.Err(err))
			return nil, err
		}

		log.Debug("user found by email")
		return user, nil
	}

	func (r *userRepository) FindByProviderID(ctx context.Context, provider domain.Provider, providerID string) (*domain.User, error) {
		log := logger.FromContext(ctx).With(
			logger.Layer("repository"),
			logger.Operation("find_user_by_provider"),
			logger.DBTable("users"),
		)

		log.Debug("querying user by provider id")

		query := `
			SELECT id, email, password_hash, full_name, provider, provider_id, is_verified, created_at, updated_at
			FROM users WHERE provider = $1 AND provider_id = $2
		`

		row := r.db.QueryRow(ctx, query, provider, providerID)
		user, err := scanUser(row)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Warn("user not found by provider id")
				return nil, domain.ErrUserNotFound
			}
			log.Error("failed to query user by provider", logger.Err(err))
			return nil, err
		}

		log.Debug("user found by provider id", logger.UserID(user.ID.String()))
		return user, nil
	}

	func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
		log := logger.FromContext(ctx).With(
			logger.Layer("repository"),
			logger.Operation("update_user"),
			logger.DBTable("users"),
			logger.UserID(user.ID.String()),
		)

		log.Debug("updating user")

		query := `
			UPDATE users
			SET email = $1, password_hash = $2, full_name = $3, is_verified = $4, updated_at = NOW()
			WHERE id = $5
			RETURNING id, email, password_hash, full_name, provider, provider_id, is_verified, created_at, updated_at
		`

		row := r.db.QueryRow(ctx, query,
			user.Email,
			user.PasswordHash,
			user.FirstName,
			user.LastName,
			user.DocumentType,
			user.DocumentNumber,
			user.PhoneNumber,
			user.IsVerified,
			user.ID,
		)

		updated, err := scanUser(row)
		if err != nil {
			log.Error("failed to update user", logger.Err(err))
			return nil, err
		}

		log.Info("user updated")
		return updated, nil
	}
*/
func scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.DocumentType,
		&u.DocumentNumber,
		&u.PhoneNumber,
		&u.Provider,
		&u.ProviderID,
		&u.IsVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
