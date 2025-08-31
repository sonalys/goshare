package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/sonalys/goshare/pkg/slog"
)

type Server struct {
	*http.Server
}

func setupServer(cfg Config, handler http.Handler) *Server {
	httpServer := http.Server{
		Addr:        cfg.AddrPort,
		ReadTimeout: cfg.ReadTimeout,
		Handler:     handler,
	}

	return &Server{
		Server: &httpServer,
	}
}

func (s *Server) ServeHTTP(ctx context.Context) {
	slog.Info(ctx, "http server listening", slog.WithString("addr", s.Addr))
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Panic(ctx, "listening")
	}
}
