package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/internal/pkg/logger"
	"github.com/sonalys/goshare/internal/pkg/otel"
)

var version string = "dev"

func init() {
	slog.SetDefault(logger.NewLogger())
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := Config{
		AddrPort:    ":8080",
		ReadTimeout: 10 * time.Second,
	}

	slog.Info("starting server", slog.String("version", version), slog.String("service_name", cfg.ServiceName))

	otelShutdown, err := otel.Initialize(ctx)
	if err != nil {
		slog.Error("failed to initialize otel", slog.Any("error", err))
		return
	}

	api := api.New(api.Dependencies{})
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

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg Config, handler http.Handler) *Server {
	httpServer := http.Server{
		Addr:        cfg.AddrPort,
		ReadTimeout: cfg.ReadTimeout,
		Handler:     handler,
	}

	return &Server{
		httpServer: &httpServer,
	}
}

func (s *Server) ServeHTTP(ctx context.Context) {
	slog.Info("http server listening", slog.Any("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.Any("error", err))
	}
}
