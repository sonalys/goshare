package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type Client struct {
	connPool *pgxpool.Pool
}

func NewClient(ctx context.Context, connStr string) (*Client, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connStr: %w", err)
	}

	cfg.ConnConfig.Tracer = tracer{}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &Client{
		connPool: dbpool,
	}, nil
}

func (c *Client) Shutdown() {
	c.connPool.Close()
}

func (c *Client) queries() *sqlc.Queries {
	return sqlc.New(c.connPool)
}

func (c *Client) transaction(ctx context.Context, f func(tx pgx.Tx) error) error {
	tx, err := c.connPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to rollback transaction")
		}
	}()

	err = f(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
