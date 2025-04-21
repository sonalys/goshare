package middlewares

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"go.opentelemetry.io/otel/trace"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			slog.ErrorContext(r.Context(), "panic recovered", slog.Any("error", rec))
			w.WriteHeader(http.StatusInternalServerError)

			var traceID trace.TraceID
			if span := trace.SpanFromContext(r.Context()); span != nil {
				traceID = span.SpanContext().TraceID()
			}

			resp := &handlers.ErrorResponse{
				TraceID: uuid.UUID(traceID),
				Errors: []handlers.Error{
					{
						Code:    handlers.ErrorCodeInternalError,
						Message: "internal server error",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				slog.ErrorContext(r.Context(), "failed to encode error response", slog.Any("error", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
