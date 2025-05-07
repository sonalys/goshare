package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type connection interface {
	transaction(ctx context.Context, f func(q *sqlc.Queries) error) error
	queries() *sqlc.Queries
}

type Postgres struct {
	*conn[*pgxpool.Pool]
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

	return &Postgres{
		conn: &conn[*pgxpool.Pool]{
			conn: dbpool,
		},
	}, nil
}
