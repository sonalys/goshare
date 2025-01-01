package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/queries"
)

type Client struct {
	connPool *pgxpool.Pool
}

func NewClient(ctx context.Context, connStr string) (*Client, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connStr: %w", err)
	}

	cfg.ConnConfig.Tracer = queryTracer{}

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

func (c *Client) queries() *queries.Queries {
	return queries.New(c.connPool)
}
