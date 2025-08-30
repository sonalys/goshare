package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ListMembersRequest struct {
	Actor    domain.ID
	LedgerID domain.ID
}

func (c *ledgerController) MembersList(ctx context.Context, req ListMembersRequest) (map[domain.ID]*domain.LedgerMember, error) {
	ctx, span := c.tracer.Start(ctx, "membersList",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	ledger, err := c.db.Ledger().Get(ctx, req.LedgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "listing members balance", err)
	}

	if !ledger.CanView(req.Actor) {
		return nil, slog.ErrorReturn(ctx, "authorizing ledger view", application.ErrUnauthorized)
	}

	slog.Info(ctx, "ledger members listed")

	return ledger.Members, nil
}
