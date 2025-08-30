package usercontroller

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/kset"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AddMembersRequest struct {
	ActorID  domain.ID
	LedgerID domain.ID
	Emails   []string
}

// TODO(invitations): Here it's a simplification of the user membership process.
// We can always invert the flow and create invitation links, so the users click themselves
// We can also send invites through the system and they accept the invite through the API.
func (c *ledgerController) MembersAdd(ctx context.Context, req AddMembersRequest) error {
	ctx, span := c.tracer.Start(ctx, "membersAdd",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	slog.Debug(ctx, "adding ledger member", slog.With("req", req))

	transaction := func(db application.Repositories) error {
		ledger, err := db.Ledger().Get(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("getting ledger: %w", err)
		}

		if !ledger.CanManageMembers(req.ActorID) {
			return fmt.Errorf("authorizing member management: %w", application.ErrUnauthorized)
		}

		users, err := db.User().ListByEmail(ctx, req.Emails)
		if err != nil {
			return fmt.Errorf("finding new members: %w", err)
		}

		newMemberIDs := kset.Select(func(u domain.User) domain.ID { return u.ID }, users...)

		if err = ledger.AddMember(req.ActorID, newMemberIDs...); err != nil {
			return fmt.Errorf("adding members: %w", err)
		}

		if err := db.Ledger().Update(ctx, ledger); err != nil {
			return fmt.Errorf("saving ledger: %w", err)
		}

		return nil
	}

	if err := c.db.Transaction(ctx, transaction); err != nil {
		return slog.ErrorReturn(ctx, "commiting transaction", err)
	}

	slog.Info(ctx, "added new members to ledger")

	return nil
}
