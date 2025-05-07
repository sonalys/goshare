package main

import (
	"context"
	"net/http"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
)

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
	slog.Info(ctx, "http server listening", slog.WithString("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Panic(ctx, "listening")
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
