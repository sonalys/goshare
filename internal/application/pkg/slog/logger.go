package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogotel "github.com/remychantenay/slog-otel"
)

const (
	LevelDebug slog.Level = -4 + iota*4
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
)

func Init(level slog.Level) {
	slog.SetLogLoggerLevel(level)
	internalHandler := &fieldFormatterHandler{
		Next: slogotel.OtelHandler{
			Next: tint.NewHandler(os.Stdout, &tint.Options{
				Level: level,
			}),
		},
	}
	logger := slog.New(internalHandler)
	slog.SetDefault(logger)
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
