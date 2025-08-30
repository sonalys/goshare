package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/pkg/slog"
)

func NewHandler(client server.Handler, repositories *repos, serviceName string) http.Handler {
	securityMiddleware := middlewares.NewSecurityHandler(repositories.JWTRepository)

	handler, _ := server.NewServer(client, securityMiddleware,
		server.WithPathPrefix("/api/v1"),
		server.WithErrorHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
			resp := client.NewError(ctx, err)

			w.WriteHeader(resp.StatusCode)

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				slog.Error(ctx, "failed to encode error response", err)
			}
		}),
	)

	return middlewares.Wrap(handler,
		middlewares.Recoverer,
		middlewares.LogRequests,
	)
}
