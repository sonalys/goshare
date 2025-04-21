package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/sonalys/goshare/cmd/server/api"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/cmd/server/middlewares"
	"github.com/sonalys/goshare/internal/pkg/otel"
)

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
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func InitializeHandler(client *api.API, repositories *repositories, serviceName string) http.Handler {
	securityHandler := api.NewSecurityHandler(repositories.JWTRepository)

	handler, _ := handlers.NewServer(client, securityHandler,
		handlers.WithPathPrefix("/api/v1"),
	)

	return middlewares.Wrap(handler,
		recoverMiddleware,
		middlewares.Instrument(serviceName, middlewares.MakeRouteFinder(handler), &otel.Provider{}),
	)
}
