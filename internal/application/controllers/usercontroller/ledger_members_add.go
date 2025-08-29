package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/kset"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AddMembersRequest struct {
	Actor    domain.ID
	LedgerID domain.ID
	Emails   []string
}

// TODO(invitations): Here it's a simplification of the user membership process.
// We can always invert the flow and create invitation links, so the users click themselves
// We can also send invites through the system and they accept the invite through the API.
func (c *ledgerController) AddMembers(ctx context.Context, req AddMembersRequest) error {
	ctx, span := c.tracer.Start(ctx, "addMembers",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	slog.Debug(ctx, "adding ledger member", slog.With("req", req))

	transaction := func(db application.Repositories) error {
		users, err := db.User().ListByEmail(ctx, req.Emails)
		if err != nil {
			return slog.ErrorReturn(ctx, "getting users", err)
		}

		ledger, err := db.Ledger().Find(ctx, req.LedgerID)
		if err != nil {
			return slog.ErrorReturn(ctx, "finding ledger", err)
		}

		err = ledger.AddMember(req.Actor, kset.Select(func(u domain.User) domain.ID { return u.ID }, users...)...)
		if err != nil {
			return slog.ErrorReturn(ctx, "adding members", err)
		}

		if err := db.Ledger().Update(ctx, ledger); err != nil {
			return err
		}

		return nil
	}

	switch err := c.db.Transaction(ctx, transaction); {
	case err == nil:
		slog.Info(ctx, "added users to ledger")
		return nil
	default:
		return slog.ErrorReturn(ctx, "adding ledger member", err)
	}
}
