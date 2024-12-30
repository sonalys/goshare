package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type (
	Secrets struct {
	}
)

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, code int, resp handlers.ErrorResponse) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/problem+json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := handlers.ErrorResponse{
		Errors: []handlers.Error{
			newError(r, handlers.NotFound, "url not found"),
		},
	}

	writeErrorResponse(ctx, w, http.StatusNotFound, resp)
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered", slog.Any("error", r))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func InitializeHandler(api *api.API, serviceName string) http.Handler {
	strictHandlerOptions := handlers.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		ResponseErrorHandlerFunc: responseErrorHandler,
	}
	strictHandler := handlers.NewStrictHandlerWithOptions(api, nil, strictHandlerOptions)

	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)

	otelMux := &OTELMux{mux}

	handlerOptions := handlers.StdHTTPServerOptions{
		BaseRouter:       otelMux,
		BaseURL:          "/api/v1",
		ErrorHandlerFunc: responseErrorHandler,
		Middlewares: []handlers.MiddlewareFunc{
			recoverMiddleware,
		},
	}
	handler := handlers.HandlerWithOptions(strictHandler, handlerOptions)

	// Wrap the handler with OpenTelemetry propagation.
	otelHandler := otelhttp.NewHandler(handler, "/", otelhttp.WithServerName(serviceName))

	return otelHandler
}
