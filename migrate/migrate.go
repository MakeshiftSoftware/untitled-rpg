package migrate

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/pkger"
	"github.com/rs/zerolog/log"
)

// Migrate checks migration status and performs database migrations.
func Migrate(db *sqlx.DB) {
	log.Info().Msg("Starting database migration")

	dbInstance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migrate database instance")
	}

	srcInstance, err := httpfs.New(pkger.Dir("/migrate/migrations"), "")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migrate source instance")
	}

	ms, err := migrate.NewWithInstance("httpfs", srcInstance, "postgres", dbInstance)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migration service")
	}

	version, _, err := ms.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			log.Info().Msg("No migrations have been applied yet")
		} else {
			log.Fatal().Err(err).Msg("Failed to get migration version")
		}
	} else {
		log.Info().Uint("migrationVersion", version).Send()
	}

	if err := ms.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("No new migrations to apply")
		} else {
			log.Fatal().Err(err).Msg("Failed to apply migrations")
		}
	}

	log.Info().Msg("Database migration complete")
}
