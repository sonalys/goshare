package main

import (
	"context"
	"log/slog"
	"net/http"
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
