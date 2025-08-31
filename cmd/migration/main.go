package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sonalys/goshare/internal/infrastructure/postgres/migrations"
	"github.com/sonalys/goshare/pkg/secrets"
	"github.com/sonalys/goshare/pkg/slog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	slog.Init(slog.LevelDebug)
	slog.Info(ctx, "starting migration")

	if err := migrations.MigrateUp(ctx, secrets.LoadSecrets().PostgresConn); err != nil {
		slog.Panic(ctx, "could not migrate up", slog.WithError(err))
	}
}
