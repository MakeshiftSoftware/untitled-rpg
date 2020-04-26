package migrate

import (
	"untitled_rpg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/pkger"
)

// Migrate checks database migration status and performs database migrations.
func Migrate(logger logger.Logger, db *sqlx.DB) {
	logger.Info().Msg("Starting database migration")

	dbInstance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create migrate database instance")
	}

	srcInstance, err := httpfs.New(pkger.Dir("/migrate/migrations"), "")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create migrate source instance")
	}

	ms, err := migrate.NewWithInstance("httpfs", srcInstance, "postgres", dbInstance)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create migration service")
	}

	version, _, err := ms.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			logger.Info().Msg("No migrations have been applied yet")
		} else {
			logger.Fatal().Err(err).Msg("Failed to get migration version")
		}
	} else {
		logger.Info().Uint("migrationVersion", version).Send()
	}

	if err := ms.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info().Msg("No new migrations to apply")
		} else {
			logger.Fatal().Err(err).Msg("Failed to apply migrations")
		}
	}

	logger.Info().Msg("Database migration complete")
}
