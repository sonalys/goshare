package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

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

func generateStructuredStack() string {
	type StackEntry struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     int    `json:"line"`
	}

	programCounters := make([]uintptr, 8)
	n := runtime.Callers(2, programCounters)

	frames := runtime.CallersFrames(programCounters)
	stackEntries := make([]StackEntry, 0, n)

	for {
		frame, more := frames.Next()
		stackEntries = append(stackEntries, StackEntry{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}

	stackJSON, _ := json.Marshal(stackEntries[2:])
	return string(stackJSON)
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", slog.Any("error", rec), slog.String("stack", generateStructuredStack()))
				responseErrorHandler(w, r, fmt.Errorf("panic recovered: %v", rec))
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
		ErrorHandlerFunc: requestErrorHandler,
		Middlewares: []handlers.MiddlewareFunc{
			recoverMiddleware,
		},
	}
	handler := handlers.HandlerWithOptions(strictHandler, handlerOptions)

	// Wrap the handler with OpenTelemetry propagation.
	otelHandler := otelhttp.NewHandler(handler, "HTTP",
		otelhttp.WithServerName(serviceName),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s %s", operation, r.Method, r.URL.Path)
		}),
	)

	return otelHandler
}
