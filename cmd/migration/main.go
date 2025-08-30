package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	if err := migrations.MigrateUp(ctx, secrets.LoadSecrets().PostgresConn); err != nil {
		slog.Panic(ctx, "could not migrate up", slog.WithError(err))
	}
}
