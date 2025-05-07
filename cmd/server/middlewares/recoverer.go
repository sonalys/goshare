package middlewares

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"go.opentelemetry.io/otel/trace"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			stackTrace := string(debug.Stack())
			slog.Error(r.Context(), "panic recovered", nil, slog.WithAny("r", rec), slog.WithString("stack_trace", stackTrace))
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
				slog.Error(r.Context(), "failed to encode error response", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
