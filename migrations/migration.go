package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var FS embed.FS

func RunMigrations(db *sql.DB) error {
	slog.Info("Initializing database migration driver...")

	d, err := iofs.New(FS, ".")
	if err != nil {
		slog.Error("Failed to create iofs migration source", "error", err)
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Failed to create postgres driver instance for migrations", "error", err)

		return err
	}
	m, err := migrate.NewWithInstance(
		"iofs", d,
		"postgres", driver,
	)
	if err != nil {
		slog.Error("Failed to initialize migration instance", "error", err)

		return err
	}

	slog.Info("Running database migrations...")

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Database schema is already up to date. No new migrations applied.")

			return nil
		}
		slog.Error("Database migration execution failed", "error", err)

		return err
	}
	slog.Info("Database migrations applied successfully!")

	return nil
}
