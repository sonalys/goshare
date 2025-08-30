package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/pkg/slog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *ledgerController) ListByUser(ctx context.Context, actorID domain.ID) ([]domain.Ledger, error) {
	ctx, span := c.tracer.Start(ctx, "listbyUser",
		trace.WithAttributes(
			attribute.Stringer("actor_id", actorID),
		),
	)
	defer span.End()

	ledgers, err := c.db.Ledger().ListByUser(ctx, actorID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "listing ledgers", err)
	}

	return ledgers, nil
}
