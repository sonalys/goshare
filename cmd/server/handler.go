package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"path"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
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

	writeErrorResponse(ctx, w, http.StatusNotFound, newErrorResponse(r, []handlers.Error{
		newError(handlers.NotFound, "url not found"),
	}))
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

func errorRespValidator(handler nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		resp, err := handler(ctx, w, r, request)
		if err != nil {
			return nil, err
		}

		type alias = struct {
			handlers.ErrorResponseJSONResponse
		}

		if errorResp, ok := resp.(alias); ok {
			errorResp.Url = r.URL.Path
			errorResp.TraceId = types.UUID(trace.SpanContextFromContext(r.Context()).TraceID())
			return errorResp, nil
		}

		return resp, nil
	}
}

func InitializeHandler(api *api.API, serviceName string) http.Handler {
	strictHandlerOptions := handlers.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		ResponseErrorHandlerFunc: responseErrorHandler,
	}

	strictMiddlewares := []handlers.StrictMiddlewareFunc{
		errorRespValidator,
	}
	strictHandler := handlers.NewStrictHandlerWithOptions(api, strictMiddlewares, strictHandlerOptions)

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
	otelHandler := otelhttp.NewHandler(handler, "/",
		otelhttp.WithServerName(serviceName),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return path.Join(operation, r.URL.Path)
		}),
	)

	return otelHandler
}
