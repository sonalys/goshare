package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/queries"
	"github.com/sonalys/goshare/internal/pkg/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	connPool *pgxpool.Pool
}

type queryTracer struct{}

// TraceQueryEnd implements pgx.QueryTracer.
func (t queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	if data.Err != nil {
		span.SetStatus(codes.Error, data.Err.Error())
	}
	span.End()
}

// TraceQueryStart implements pgx.QueryTracer.
func (t queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx, _ = otel.Tracer.Start(ctx, "postgres.query")
	return ctx
}

var _ pgx.QueryTracer = queryTracer{}

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
