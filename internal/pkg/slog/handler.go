package slog

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type fieldFormatterHandler struct {
	Next slog.Handler
}

func (h *fieldFormatterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Next.Enabled(ctx, level)
}

func (h *fieldFormatterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.Next.WithAttrs(attrs)
}

func (h *fieldFormatterHandler) WithGroup(name string) slog.Handler {
	return h.Next.WithGroup(name)
}

func (h *fieldFormatterHandler) Handle(ctx context.Context, record slog.Record) error {
	if ctx == nil {
		return h.Next.Handle(ctx, record)
	}

	fields := getFields(ctx)
	record.AddAttrs(fields...)
	if record.Level >= slog.LevelError {
		if span := trace.SpanFromContext(ctx); span != nil {
			span.SetStatus(codes.Error, record.Message)
		}
	}

	return h.Next.Handle(ctx, record)
}
