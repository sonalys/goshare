package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/pkg/slog"
	v1 "github.com/sonalys/goshare/pkg/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetLedgerRequest struct {
	ActorID  domain.ID
	LedgerID domain.ID
}

func (c *ledgerController) Get(ctx context.Context, req GetLedgerRequest) (*domain.Ledger, error) {
	ctx, span := c.tracer.Start(ctx, "get",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	ledger, err := c.db.Ledger().Get(ctx, req.LedgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "listing ledgers", err)
	}

	if !ledger.CanView(req.ActorID) {
		return nil, slog.ErrorReturn(ctx, "authorizing ledger view", v1.ErrForbidden)
	}

	return ledger, nil
}
