package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/application/pkg/otel"
	"github.com/sonalys/goshare/internal/application/pkg/secrets"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
)

var version string = "dev"

func init() {
	slog.Init(slog.LevelDebug)
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

	telemetryShutdown, err := otel.Initialize(ctx, cfg.TelemetryEndpoint, version)
	if err != nil {
		slog.Panic(ctx, "starting telemetry", slog.WithError(err))
	}

	secrets := secrets.LoadSecrets()

	infrastructure := loadInfrastructure(ctx, secrets)
	repositories := loadRepositories(secrets, infrastructure)

	userController := usercontroller.New(usercontroller.Dependencies{
		Database: repositories.Database,
	})

	identityController := identitycontroller.New(identitycontroller.Dependencies{
		Database:        repositories.Database,
		IdentityEncoder: repositories.JWTRepository,
	})

	api := api.New(
		identityController,
		userController,
	)
	handler := NewHandler(api, repositories, cfg.ServiceName)

	server := NewServer(cfg, handler)

	go server.ServeHTTP(ctx)

	<-ctx.Done()

	shutdown(
		server.Shutdown,
		telemetryShutdown,
	)
}

func shutdown(fns ...func(context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	wait := make(chan struct{})

	for _, fn := range fns {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) {
				slog.Error(ctx, "shutting down service", err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(wait)
	}()

	select {
	case <-ctx.Done():
		slog.Error(ctx, "could not stop all services", nil)
		syscall.Exit(1)
	case <-wait:
		slog.Info(ctx, "all services stopped gracefully")
		syscall.Exit(0)
	}
}
