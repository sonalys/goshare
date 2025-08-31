package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/pkg/otel"
	"github.com/sonalys/goshare/pkg/slog"
)

func setupHandler(ctx context.Context, securityHandler server.SecurityHandler, client server.Handler) http.Handler {
	handler, err := server.NewServer(client, securityHandler,
		server.WithPathPrefix("/api/v1"),
		server.WithTracerProvider(otel.Provider.TracerProvider()),
		server.WithMiddleware(
			middlewares.Recoverer,
			middlewares.Logger,
		),
		server.WithErrorHandler(errorHandler(client)),
	)
	if err != nil {
		slog.Panic(ctx, "creating http api handler", slog.WithError(err))
	}

	return handler
}

func errorHandler(client server.Handler) func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		resp := client.NewError(ctx, err)

		w.WriteHeader(resp.StatusCode)

		w.Header().Set("Content-Type", "application/json")
		//nolint:musttag // generated structure.
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(ctx, "failed to encode error response", err)
		}
	}
}
