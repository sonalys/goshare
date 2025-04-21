package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogotel "github.com/remychantenay/slog-otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type fieldFormatterHandler struct {
	slog.Handler
}

func NewLogger() *slog.Logger {
	handler := tint.NewHandler(os.Stdout, nil)
	otelHandler := slogotel.OtelHandler{
		Next: handler,
	}
	internalHandler := &fieldFormatterHandler{otelHandler}
	return slog.New(internalHandler)
}

func (h *fieldFormatterHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= slog.LevelError {
		if span := trace.SpanFromContext(ctx); span != nil {
			span.SetStatus(codes.Error, record.Message)
		}
	}
	return h.Handler.Handle(ctx, record)
}
