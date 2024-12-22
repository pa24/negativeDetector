package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

// RunMigrations запускает миграции с использованием goose
func RunMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	migrationsDir := "../internal/database/migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}
