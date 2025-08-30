package slog

import (
	"fmt"
	"log/slog"
	"time"
)

func WithError(err error) slog.Attr {
	if err == nil {
		return slog.Any("error", nil)
	}
	return slog.Any("error", err.Error())
}

func With(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func WithStringer(key string, str fmt.Stringer) slog.Attr {
	return slog.String(key, str.String())
}

func WithString(key string, value string) slog.Attr {
	return slog.String(key, value)
}

func WithInt(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func WithDuration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

func WithUint64(key string, value uint64) slog.Attr {
	return slog.Uint64(key, value)
}

func WithBool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}
