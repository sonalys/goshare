package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type (
	Secrets struct {
	}
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.WriteHeader(http.StatusNotFound)

	resp := handlers.ErrorResponse{
		Errors: []handlers.Error{
			{
				TraceId: types.UUID(trace.SpanContextFromContext(ctx).TraceID()),
				Code:    handlers.NotFound,
				Message: "url not found",
				Url:     r.URL.Path,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
	}
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
	}
	handler := handlers.HandlerWithOptions(strictHandler, handlerOptions)

	// Wrap the handler with OpenTelemetry propagation.
	otelHandler := otelhttp.NewHandler(handler, "/", otelhttp.WithServerName(serviceName))

	return otelHandler
}
