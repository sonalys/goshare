package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type connection interface {
	transaction(ctx context.Context, f func(q *sqlc.Queries) error) error
	queries() *sqlc.Queries
}

type Postgres struct {
	*queries[*pgxpool.Pool]
}

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

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if dbpool.Ping(ctx) == nil {
			break
		}

		slog.Info(ctx, "waiting for postgres connection")
		time.Sleep(time.Second)
	}

	return &Postgres{
		&queries[*pgxpool.Pool]{
			conn[*pgxpool.Pool]{
				conn: dbpool,
			},
		},
	}, nil
}
