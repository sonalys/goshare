package ledgers

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	AddMembersRequest struct {
		UserID   v1.ID
		LedgerID v1.ID
		Emails   []string
	}
)

func (r *AddMembersRequest) Validate() error {
	var errs v1.FormError

	if r.UserID.IsEmpty() {
		errs = append(errs, v1.NewRequiredFieldError("user_id"))
	}

	if r.LedgerID.IsEmpty() {
		errs = append(errs, v1.NewRequiredFieldError("ledger_id"))
	}

	r.Emails = slices.Compact(r.Emails)
	switch lenEmails := len(r.Emails); {
	case lenEmails == 0:
		errs = append(errs, v1.NewRequiredFieldError("emails"))
	case lenEmails > v1.LedgerMaxUsers-1:
		errs = append(errs, v1.NewFieldLengthError("emails", 1, v1.LedgerMaxUsers))
	}

	return errs.Validate()
}

// TODO(invitations): Here it's a simplification of the user membership process.
// We can always invert the flow and create invitation links, so the users click themselves
// We can also send invites through the system and they accept the invite through the API.
func (c *Controller) AddParticipants(ctx context.Context, req AddMembersRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.AddMembers")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("user_id", req.UserID),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	users, err := c.userRepository.ListByEmail(ctx, req.Emails)
	if err != nil {
		return slog.ErrorReturn(ctx, "failed to get users by email", err)
	}

	ids := make([]v1.ID, 0, len(users))
	for _, user := range users {
		if user.ID == req.UserID {
			continue
		}
		ids = append(ids, user.ID)
	}

	err = c.ledgerRepository.AddParticipants(ctx, req.LedgerID, func(ledger *v1.Ledger) error {
		if ledger.CreatedBy != req.UserID {
			return fmt.Errorf("user %s is not the owner of the ledger %s", req.UserID, ledger.ID)
		}

		ledger.AddParticipants(req.UserID, ids...)

		if len(ledger.Participants) >= v1.LedgerMaxUsers {
			return v1.ErrLedgerMaxUsers
		}

		return nil
	})
	switch {
	case err == nil:
		slog.Info(ctx, "added users to ledger")
		return nil
	case errors.Is(err, v1.ErrNotFound):
		return v1.FieldError{
			Field: "ledger_id",
			Cause: v1.ErrNotFound,
		}
	default:
		return slog.ErrorReturn(ctx, "failed to add users to ledger", err)
	}
}

func (c *Controller) GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetBalances")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("ledger_id", ledgerID),
	)

	participants, err := c.ledgerRepository.GetParticipants(ctx, ledgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger participants balances", err)
	}

	slog.Info(ctx, "ledger participants balances retrieved")

	return participants, nil
}
