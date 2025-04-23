package main

import (
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/pkg/logger"
	"github.com/sonalys/goshare/internal/pkg/secrets"
)

func init() {
	logger.InitializeLogger()
}

func main() {
	slog.Info("starting migration")

	driver, err := iofs.New(postgres.MigrationsFS, "migrations")
	if err != nil {
		panic(err)
	}

	slog.Info("migration files loaded")

	secrets := secrets.LoadSecrets()

	m, err := migrate.NewWithSourceInstance("iofs", driver, secrets.PostgresConn)
	if err != nil {
		panic(err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		panic(err)
	}

	slog.Info("migrate driver loaded", slog.Uint64("currentVersion", uint64(version)), slog.Bool("isDirty", dirty))

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("no changes to migrate")
			return
		}
		panic(err)
	}

	version, _, err = m.Version()
	if err != nil {
		panic(err)
	}

	slog.Info("migrated up", slog.Uint64("currentVersion", uint64(version)))
}
