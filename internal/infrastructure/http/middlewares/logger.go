package middlewares

import (
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/pkg/slog"
	"go.opentelemetry.io/otel/trace"
)

func Logger(req middleware.Request, next middleware.Next) (resp middleware.Response, err error) {
	ctx := req.Context

	fields := []any{
		slog.WithString("operation_id", req.OperationID),
	}

	slog.Info(ctx, "request received", fields...)

	resp, err = next(req)
	switch err := err.(type) {
	case nil:
		if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
			fields = append(fields,
				slog.WithInt("status_code", tresp.GetStatusCode()),
			)
		}
	case *server.ErrorResponseStatusCode:
		var traceID trace.TraceID
		if span := trace.SpanFromContext(ctx); span != nil {
			traceID = span.SpanContext().TraceID()
		}
		err.Response.TraceID = uuid.UUID(traceID)
	}

	slog.Info(ctx, "request completed", fields...)

	return
}
