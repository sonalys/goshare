package middlewares

import (
	"log/slog"
	"net/http"
)

// LogRequests logs incoming requests using context logger.
func LogRequests(find RouteFinder) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var opID, opName string
			if route, ok := find(r.Method, r.URL); ok {
				opName = route.Name()
				opID = route.OperationID()
			}
			slog.InfoContext(r.Context(), "request received",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				opID,
				opName,
			)
			next.ServeHTTP(w, r)
		})
	}
}
