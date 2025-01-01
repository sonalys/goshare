package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/lmittmann/tint"
	slogotel "github.com/remychantenay/slog-otel"
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

func generateStack(pc uintptr) string {
	type StackEntry struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     string `json:"line"`
	}

	frames := runtime.CallersFrames([]uintptr{pc})
	stack := make([]StackEntry, 0, 1)

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

	stackJSON, _ := json.Marshal(stack)

	return string(stackJSON)
}

func (h *fieldFormatterHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= slog.LevelError {
		var hasStack bool
		record.Attrs(func(a slog.Attr) bool {
			if a.Key == "stack" {
				hasStack = true
				return false
			}
			return true
		})

		if !hasStack {
			record.AddAttrs(slog.Any("stack", generateStack(record.PC)))
		}
	}

	return h.Handler.Handle(ctx, record)
}
