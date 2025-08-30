package slog

import (
	"context"
	"log/slog"
)

type contextKey string

const privateKey = contextKey("log-metadata")

func getFields(ctx context.Context) []slog.Attr {
	fields, _ := ctx.Value(privateKey).([]slog.Attr)
	return fields
}

func Context(ctx context.Context, fields ...slog.Attr) context.Context {
	fields = append(fields, getFields(ctx)...)
	return context.WithValue(ctx, privateKey, fields)
}
