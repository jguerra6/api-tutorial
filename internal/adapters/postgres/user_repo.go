package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/jguerra6/api-tutorial/internal/domain"
	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
	"github.com/jguerra6/api-tutorial/internal/ports"
)

var _ ports.UserRepository = &repository{}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) ports.UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Insert(ctx context.Context, u *domain.User) error {

	logger := ctxutils.Logger(ctx).With().
		Str("component", "repo.user").
		Str("op", "insert").
		Str("entity", "user").
		Logger()

	row, err := userFromDomain(u)
	if err != nil {
		logger.Error().Err(err)
		return err
	}

	tx, err := r.db.Beginx()
	if err != nil {
		logger.Error().Err(err)
		return ErrSQLBegin
	}

	const stmt = `INSERT INTO users (id, email, role, phone_number, display_name, is_active, is_verified, created_at, updated_at) 
			VALUES (:id, :email, :role, :phone_number, :display_name, :is_active, :is_verified, :created_at, :updated_at) RETURNING id, created_at`

	_, err = tx.NamedExecContext(ctx, stmt, &row)
	if err != nil {
		logger.Error().Err(err).Msg("tx.QueryxContext error")
		if err = tx.Rollback(); err != nil {
			logger.Debug().Err(err).Msg("tx.Rollback() error")
			return ErrRollback
		}
		return ErrSQLExec
	}

	err = tx.Commit()
	if err != nil {
		logger.Debug().Err(err).Msg("tx.Commit() error")
		return ErrSQLExec
	}

	logger.Info().Msg("db insert ok")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {

	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	const stmt = `DELETE FROM users WHERE id = $1 RETURNING 1`

	tx, err := r.db.Beginx()
	if err != nil {
		return ErrSQLBegin
	}
	var deleted int
	err = tx.QueryRowxContext(ctx, stmt, uid).Scan(&deleted)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {

		if err = tx.Rollback(); err != nil {
			return ErrRollback
		}
		return ErrSQLExec

	}

	if deleted == 0 {
		err := tx.Rollback()
		if err != nil {
			return ErrRollback
		}
		return ports.NewNotFoundError(ErrDataNotFound.Error())
	}

	if err = tx.Commit(); err != nil {
		return ErrSQLExec
	}

	return nil
}
