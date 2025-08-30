package migrations

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sonalys/goshare/internal/application/pkg/slog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed *.sql
var migrationFS embed.FS

func MigrateUp(ctx context.Context, connStr string) error {
	driver, err := iofs.New(migrationFS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", driver, connStr)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	slog.Info(ctx, "migrate driver loaded", slog.WithUint64("current_version", uint64(version)), slog.WithBool("is_dirty", dirty))

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info(ctx, "no changes to migrate")
			return nil
		}
		slog.Panic(ctx, "migrating up")
	}

	version, _, err = m.Version()
	if err != nil {
		slog.Panic(ctx, "reading current version")
	}

	slog.Info(ctx, "migrated up", slog.WithUint64("currentVersion", uint64(version)))

	return nil
}
