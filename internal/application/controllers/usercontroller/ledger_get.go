package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetLedgerRequest struct {
	Actor    domain.ID
	LedgerID domain.ID
}

func (c *ledgerController) Get(ctx context.Context, req GetLedgerRequest) (*domain.Ledger, error) {
	ctx, span := c.tracer.Start(ctx, "get",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	ledger, err := c.db.Ledger().Get(ctx, req.LedgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "listing ledgers", err)
	}

	return ledger, nil
}
