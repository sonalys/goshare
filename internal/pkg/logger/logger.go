package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"

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

	if record.Level >= slog.LevelError {
		type StackEntry struct {
			Function string `json:"function"`
			File     string `json:"file"`
			Line     string `json:"line"`
		}

		frames := runtime.CallersFrames([]uintptr{record.PC})
		var stack []StackEntry
		for {
			frame, more := frames.Next()
			stack = append(stack, StackEntry{
				Function: frame.Function,
				File:     frame.File,
				Line:     fmt.Sprint(frame.Line),
			})
			if !more {
				break
			}
		}

		record.AddAttrs(slog.Any("stack", stack))
	}

	return h.Handler.Handle(ctx, record)
}
