package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	api.WriteErrorResponse(ctx, w, http.StatusNotFound, newErrorResponse(r, []handlers.Error{
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

func InitializeHandler(client *api.API, repositories *repositories, serviceName string) http.Handler {
	strictHandlerOptions := handlers.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		ResponseErrorHandlerFunc: responseErrorHandler,
	}

	strictMiddlewares := []handlers.StrictMiddlewareFunc{
		api.AuthMiddleware(repositories.JWTRepository),
		api.InjectRequestContextDataMiddleware, // Should be the last middleware, so all other middlewares can access the context data.
	}
	strictHandler := handlers.NewStrictHandlerWithOptions(client, strictMiddlewares, strictHandlerOptions)

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
	otelHandler := otelhttp.NewHandler(handler, "HTTP",
		otelhttp.WithServerName(serviceName),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s", operation, r.URL.Path)
		}),
	)

	return otelHandler
}
