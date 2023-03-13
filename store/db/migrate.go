package db

import (
	"database/sql"
	"embed"

	"github.com/fox-one/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed schema/*.sql
var fs embed.FS

// Migrate run db migration with embed schemes.
func Migrate(db *sql.DB) error {
	d, err := iofs.New(fs, "schema")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgresql", driver)
	if err != nil {
		return err
	}

	m.Log = &migrateLog{}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

type migrateLog struct{}

func (migrateLog) Printf(format string, v ...interface{}) {
	logger.L.Printf(format, v...)
}

func (migrateLog) Verbose() bool {
	return true
}
