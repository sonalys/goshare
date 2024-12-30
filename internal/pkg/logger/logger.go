package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/otel/trace"
)

type otelHandler struct {
	slog.Handler
}

func NewLogger() *slog.Logger {
	handler := tint.NewHandler(os.Stdout, nil)

	internalHandler := &otelHandler{handler}

	return slog.New(internalHandler)
}

func (h *otelHandler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanContextFromContext(ctx)

	if span.HasSpanID() {
		record.AddAttrs(slog.String("span_id", span.SpanID().String()))
	}

	if span.HasTraceID() {
		record.AddAttrs(slog.String("trace_id", span.TraceID().String()))
	}

	return h.Handler.Handle(ctx, record)
}
