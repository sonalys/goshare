package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
)

type Config struct {
	ServiceName       string
	AddrPort          string
	ReadTimeout       time.Duration
	EnableTelemetry   bool
	TelemetryEndpoint string
}

func loadConfigFromEnv(ctx context.Context) Config {
	cfg := Config{
		ServiceName:       "server",
		AddrPort:          ":8080",
		ReadTimeout:       10 * time.Second,
		EnableTelemetry:   true,
		TelemetryEndpoint: "jaeger:4317",
	}

	if serviceName, ok := os.LookupEnv("SERVICE_NAME"); ok {
		cfg.ServiceName = serviceName
	}

	if enableTelemetry, ok := os.LookupEnv("ENABLE_TELEMETRY"); ok {
		enableTelemetry, ok := strconv.ParseBool(enableTelemetry)
		if ok == nil {
			cfg.EnableTelemetry = enableTelemetry
		} else {
			slog.Warn(ctx, "failed to parse ENABLE_TELEMETRY")
		}
	}

	if addrPort, ok := os.LookupEnv("ADDR_PORT"); ok {
		cfg.AddrPort = addrPort
	}

	if readTimeout, ok := os.LookupEnv("READ_TIMEOUT"); ok {
		if d, err := time.ParseDuration(readTimeout); err == nil {
			cfg.ReadTimeout = d
		} else {
			slog.Warn(ctx, "failed to parse READ_TIMEOUT", slog.WithError(err))
		}
	}

	return cfg
}
