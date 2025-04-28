package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogotel "github.com/remychantenay/slog-otel"
)

const LevelPanic slog.Level = 12

func Init() {
	internalHandler := &fieldFormatterHandler{
		Next: slogotel.OtelHandler{
			Next: tint.NewHandler(os.Stdout, nil),
		},
	}
	slog.SetDefault(slog.New(internalHandler))
}

func Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, err error, args ...any) {
	if err != nil {
		args = append(args, WithError(err))
	}
	slog.ErrorContext(ctx, msg, args...)
}

func Panic(ctx context.Context, msg string, args ...any) {
	slog.Default().Log(ctx, LevelPanic, msg, args...)
	panic(msg)
}

func ErrorReturn(ctx context.Context, msg string, err error, args ...any) error {
	slog.ErrorContext(ctx, msg, append(args, WithError(err))...)
	return fmt.Errorf("%s: %w", msg, err)
}
