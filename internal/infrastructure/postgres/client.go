package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

type connection interface {
	transaction(ctx context.Context, f func(q connection) error) error
	queries() *sqlcgen.Queries
	readWrite() *readWriteRepository

	application.Queries
}

type Postgres struct {
	connection
}

var _ application.Database = &Postgres{}

func New(ctx context.Context, connStr string) (*Postgres, error) {
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

	return &Postgres{
		connection: &conn[*pgxpool.Pool]{
			conn: dbpool,
		},
	}, nil
}

func (c *Postgres) Transaction(ctx context.Context, f func(application.Repositories) error) error {
	return c.transaction(ctx, func(q connection) error {
		return f(c.readWrite())
	})
}

func wait(ctx context.Context, conn *pgxpool.Pool) error {
	for {
		if conn.Ping(ctx) == nil {
			return nil
		}

		slog.Info(ctx, "waiting for postgres connection")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}
