package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sonalys/goshare/internal/application/pkg/secrets"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/migrations"
)

func init() {
	slog.Init(slog.LevelDebug)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	slog.Info(ctx, "starting migration")

	driver, err := iofs.New(migrations.MigrationsFS, ".")
	if err != nil {
		panic(err)
	}

	slog.Info(ctx, "migration files loaded")

	secrets := secrets.LoadSecrets()

	m, err := migrate.NewWithSourceInstance("iofs", driver, secrets.PostgresConn)
	if err != nil {
		panic(err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		panic(err)
	}

	slog.Info(ctx, "migrate driver loaded", slog.WithUint64("current_version", uint64(version)), slog.WithBool("is_dirty", dirty))

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info(ctx, "no changes to migrate")
			return
		}
		slog.Panic(ctx, "migrating up")
	}

	version, _, err = m.Version()
	if err != nil {
		slog.Panic(ctx, "reading current version")
	}

	slog.Info(ctx, "migrated up", slog.WithUint64("currentVersion", uint64(version)))
}
