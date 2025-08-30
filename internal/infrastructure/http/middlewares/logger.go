package middlewares

import (
	"net/http"
	"time"

	"github.com/sonalys/goshare/pkg/slog"
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

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		slog.Info(ctx, "request received",
			slog.WithString("method", r.Method),
			slog.WithString("url", r.URL.String()),
			slog.WithString("remote_addr", r.RemoteAddr),
		)

		rw := wrapResponseWriter(w)

		t1 := time.Now()

		defer func() {
			status := rw.Status()

			ctx = slog.Context(ctx,
				slog.WithString("method", r.Method),
				slog.WithString("url", r.URL.String()),
				slog.WithString("remote_addr", r.RemoteAddr),
				slog.WithInt("status", rw.Status()),
				slog.WithDuration("duration", time.Since(t1)),
			)

			if status >= 400 {
				slog.Error(ctx, "request failed", nil)
				return
			}

			slog.Info(ctx, "request completed")
		}()

		next.ServeHTTP(rw, r)
	})
}
