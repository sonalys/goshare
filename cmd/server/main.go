package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/internal/pkg/logger"
	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/secrets"
)

var version string = "dev"

func init() {
	slog.SetDefault(logger.NewLogger())
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := loadConfigFromEnv()

	slog.Info("starting server", slog.String("version", version), slog.String("service_name", cfg.ServiceName))

	otelShutdown, err := otel.Initialize(ctx, cfg.TelemetryEndpoint)
	if err != nil {
		slog.Error("failed to initialize otel", slog.Any("error", err))
		return
	}

	secrets := secrets.LoadSecrets()

	infrastructure := loadInfrastructure(ctx, secrets)
	repositories := loadRepositories(infrastructure)
	controllers := loadControllers(repositories)

	api := api.New(api.Dependencies{
		ParticipantRegister: controllers.participantController,
	})
	handler := InitializeHandler(api, cfg.ServiceName)

	server := NewServer(cfg, handler)

	go server.ServeHTTP(ctx)

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	otelShutdown(ctx)

	slog.Info("shutting down server")
}
