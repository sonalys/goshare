package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/secrets"
	"github.com/sonalys/goshare/internal/pkg/slog"
)

var version string = "dev"

func init() {
	slog.Init()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := loadConfigFromEnv(ctx)

	ctx = slog.Context(ctx,
		slog.WithString("version", version),
		slog.WithString("service_name", cfg.ServiceName),
	)

	slog.Info(ctx, "starting server")

	shutdown, err := otel.Initialize(ctx, cfg.TelemetryEndpoint)
	if err != nil {
		slog.Panic(ctx, "starting telemetry", slog.WithError(err))
	}

	secrets := secrets.LoadSecrets()

	infrastructure := loadInfrastructure(ctx, secrets)
	repositories := loadRepositories(secrets, infrastructure)
	controllers := loadControllers(repositories)

	api := api.New(api.Dependencies{
		UserController:   controllers.userController,
		LedgerController: controllers.ledgerController,
	})
	handler := NewHandler(api, repositories, cfg.ServiceName)

	server := NewServer(cfg, handler)

	go server.ServeHTTP(ctx)

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	if err := shutdown(ctx); err != nil {
		slog.Error(ctx, "stopping telemetry", err)
	}

	slog.Info(ctx, "shutting down")
}
