package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/cmd/server/middlewares"
	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/slog"
)

func NewHandler(client *api.API, repositories *repositories, serviceName string) http.Handler {
	securityHandler := api.NewSecurityHandler(repositories.JWTRepository)

	handler, _ := handlers.NewServer(client, securityHandler,
		handlers.WithPathPrefix("/api/v1"),
		handlers.WithErrorHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
			resp := client.NewError(ctx, err)

			w.WriteHeader(resp.StatusCode)

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				slog.Error(ctx, "failed to encode error response", err)
			}
		}),
	)

	finder := middlewares.MakeRouteFinder(handler)

	return middlewares.Wrap(handler,
		middlewares.Instrument(serviceName, finder, otel.Provider{}),
		middlewares.Recoverer,
		middlewares.LogRequests(finder),
	)
}
