package main

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	ServiceName string
	AddrPort    string
	ReadTimeout time.Duration
}

func NewConfigFromEnv() Config {
	cfg := Config{
		ServiceName: "server",
		AddrPort:    ":8080",
		ReadTimeout: 10 * time.Second,
	}

	if serviceName, ok := os.LookupEnv("SERVICE_NAME"); ok {
		cfg.ServiceName = serviceName
	}

	if addrPort, ok := os.LookupEnv("ADDR_PORT"); ok {
		cfg.AddrPort = addrPort
	}

	if readTimeout, ok := os.LookupEnv("READ_TIMEOUT"); ok {
		if d, err := time.ParseDuration(readTimeout); err == nil {
			cfg.ReadTimeout = d
		} else {
			slog.Warn("failed to parse READ_TIMEOUT", slog.Any("error", err))
		}
	}

	return cfg
}
