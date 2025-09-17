package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	driverName = "pgx"
)

type PoolOption func(*pgxpool.Config)

func NewDb(_ context.Context, env, host, username, password, name, port string) (*sqlx.DB, error) {
	dsn := BuildDSN(env, host, username, password, name, port)

	pool, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func BuildDSN(env, host, username, password, name, port string) string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC",
		host, username, password, name, port,
	)

	if env == "production" {
		return dsn + " sslmode=verify-full sslrootcert=/app/rds-ca-bundle.pem"
	}

	return dsn + " sslmode=disable"
}
