package ledgers

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
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
func (c *Controller) AddMembers(ctx context.Context, req AddMembersRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.AddMembers")
	defer span.End()

	users, err := c.userRepository.GetByEmail(ctx, req.Emails)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to get users by email")
		return fmt.Errorf("failed to get users by email: %w", err)
	}

	var errs v1.FormError

	for i := range users {
		req.Emails = slices.DeleteFunc(req.Emails, func(email string) bool {
			return email == users[i].Email
		})

		attrs := []any{
			slog.String("user_id", users[i].ID.String()),
			slog.String("ledger_id", req.LedgerID.String()),
			slog.String("req_user_id", req.UserID.String()),
		}
		err := c.ledgerRepository.AddParticipant(ctx, req.LedgerID, req.UserID, users[i].ID)
		if err == nil {
			slog.InfoContext(ctx, "added user to ledger", attrs...)
			continue
		}

		slog.ErrorContext(ctx, "failed to add user to ledger", append(attrs,
			slog.Any("error", err),
		)...,
		)

		errs = append(errs, v1.FieldError{
			Field: fmt.Sprintf("emails.%d", i),
			Cause: err,
		})
	}

	for i := range req.Emails {
		errs = append(errs, v1.FieldError{
			Field: fmt.Sprintf("emails.%d", i),
			Cause: v1.ErrNotFound,
		})
	}

	if err := errs.Validate(); err != nil {
		span.SetStatus(codes.Error, "failed to add users to ledger")
		return err
	}

	return nil
}
