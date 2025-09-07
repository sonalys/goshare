package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/pkg/slog"
)

func New(ctx context.Context, connStr string) (Connection, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connStr: %w", err)
	}

	cfg.ConnConfig.Tracer = tracer{}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := wait(ctx, dbpool); err != nil {
		return nil, fmt.Errorf("waiting for postgres connection: %w", err)
	}

	return &conn[*pgxpool.Pool]{
		conn: dbpool,
	}, nil
}

func wait(ctx context.Context, conn *pgxpool.Pool) error {
	for {
		err := conn.Ping(ctx)
		if err == nil {
			return nil
		}

		slog.Info(ctx, "waiting for postgres connection", slog.WithError(err))

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}
