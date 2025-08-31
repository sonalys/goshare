package middlewares

import (
	"errors"

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
	if err == nil {
		if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
			fields = append(fields,
				slog.WithInt("status_code", tresp.GetStatusCode()),
			)
		}
	}

	if target := new(server.ErrorResponseStatusCode); errors.As(err, &target) {
		var traceID trace.TraceID
		if span := trace.SpanFromContext(ctx); span != nil {
			traceID = span.SpanContext().TraceID()
		}
		target.Response.TraceID = uuid.UUID(traceID)
	}

	slog.Info(ctx, "request completed", fields...)

	return
}
