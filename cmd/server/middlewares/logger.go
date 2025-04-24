package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func LogRequests(find RouteFinder) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var opID string
			if route, ok := find(r.Method, r.URL); ok {
				opID = route.OperationID()
			}
			slog.InfoContext(ctx, "request received",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("operation_id", opID),
				slog.String("remote_addr", r.RemoteAddr),
			)

			rw := wrapResponseWriter(w)

			t1 := time.Now()

			defer func() {
				status := rw.Status()

				fields := []any{
					slog.String("method", r.Method),
					slog.String("url", r.URL.String()),
					slog.String("operation_id", opID),
					slog.String("remote_addr", r.RemoteAddr),
					slog.Int("status", rw.Status()),
					slog.Duration("duration", time.Since(t1)),
				}

				if status >= 400 {
					slog.ErrorContext(ctx, "request failed", fields...)
					return
				}

				slog.InfoContext(ctx, "request completed", fields...)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
