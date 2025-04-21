package ledgers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/sonalys/goshare/internal/pkg/otel"
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

	logFields := []any{
		slog.String("user_id", req.UserID.String()),
		slog.String("ledger_id", req.LedgerID.String()),
	}

	users, err := c.userRepository.GetByEmail(ctx, req.Emails)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get users by email", append(logFields, slog.Any("error", err))...)
		return fmt.Errorf("failed to get users by email: %w", err)
	}

	ids := make([]v1.ID, 0, len(users))
	for _, user := range users {
		if user.ID == req.UserID {
			continue
		}
		ids = append(ids, user.ID)
	}

	err = c.ledgerRepository.AddParticipants(ctx, req.LedgerID, req.UserID, ids...)
	switch {
	case err == nil:
		slog.InfoContext(ctx, "added users to ledger", logFields...)
		return nil
	case errors.Is(err, v1.ErrNotFound):
		return v1.FieldError{
			Field: "ledger_id",
			Cause: v1.ErrNotFound,
		}
	default:
		slog.ErrorContext(ctx, "failed to add users to ledger", append(logFields, slog.Any("error", err))...)
		return fmt.Errorf("failed to add users to ledger: %w", err)
	}
}
