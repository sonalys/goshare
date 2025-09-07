package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/router"
	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/secrets"
	"github.com/sonalys/goshare/internal/pkg/slog"
)

var version string = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	slog.Init(slog.LevelDebug)
	cfg := loadConfigFromEnv(ctx)

	slog.Info(ctx, "starting server",
		slog.WithString("version", version),
		slog.WithString("service_name", cfg.ServiceName),
	)

	telemetryShutdown, err := otel.Initialize(ctx, cfg.TelemetryEndpoint, version)
	if err != nil {
		slog.Panic(ctx, "starting telemetry", slog.WithError(err))
	}

	secrets := secrets.LoadSecrets()

	infrastructure := loadInfrastructure(ctx, secrets)
	repositories := loadRepositories(secrets, infrastructure)

	userController := usercontroller.New(usercontroller.Dependencies{
		LocalDatabase: repositories.Database,
	})

	identityController := identitycontroller.New(identitycontroller.Dependencies{
		LocalDatabase:   repositories.Database,
		IdentityEncoder: repositories.JWTRepository,
	})

	securityHandler := middlewares.NewSecurityHandler(repositories.JWTRepository)

	router := router.New(
		securityHandler,
		identityController,
		userController,
	)
	handler := setupHandler(ctx, securityHandler, router)
	server := setupServer(cfg, handler)

	go server.ServeHTTP(ctx)

	<-ctx.Done()

	slog.Info(ctx, "shutdown signal received")

	shutdown(
		server.Shutdown,
		telemetryShutdown,
	)
}

func shutdown(fns ...func(context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var waitGroup sync.WaitGroup

	wait := make(chan struct{})

	for _, fn := range fns {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) {
				slog.Error(ctx, "shutting down service", err)
			}
		}()
	}

	go func() {
		waitGroup.Wait()
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
