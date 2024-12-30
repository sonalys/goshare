package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger() *slog.Logger {
	handler := tint.NewHandler(os.Stdout, nil)

	return slog.New(handler)
}
