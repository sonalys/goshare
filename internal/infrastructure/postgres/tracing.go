package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/pkg/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type queryTracer struct{}

var _ pgx.QueryTracer = queryTracer{}

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
