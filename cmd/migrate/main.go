package main

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jguerra6/api-tutorial/config"
	dbAdapter "github.com/jguerra6/api-tutorial/internal/adapters/postgres"
)

func main() {

	zerolog.TimeFieldFormat = time.RFC3339
	var (
		logger = log.With().Timestamp().Logger()
		ctx    = context.Background()
	)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	dsn := dbAdapter.BuildDSN("", cfg.DbHost, cfg.DbUserName, cfg.DbPassword, cfg.DbName, cfg.DbPort)

	_ = ctx
	_ = logger

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database: ")
	}

	err = db.AutoMigrate(
		dbAdapter.UserRow{},
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database: ")
	}

	log.Info().Msg("Database migration complete!")

}
