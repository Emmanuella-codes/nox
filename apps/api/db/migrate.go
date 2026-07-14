package db

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(databaseURL string) error {
	databaseURL = normalizeDatabaseURL(databaseURL)

	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, dbErr := migrator.Close()
		if sourceErr != nil {
			log.Error().Err(sourceErr).Msg("failed to close migration source")
		}
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("failed to close migration database driver")
		}
	}()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Info().Msg("database migrations applied")
	return nil
}
