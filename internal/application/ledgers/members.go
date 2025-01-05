package ledgers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
)

type (
	AddMembersRequest struct {
		UserID   uuid.UUID
		LedgerID uuid.UUID
		Emails   []string
	}
)

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

	var errList v1.FormError

	for i := range users {
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
		)...)

		errList.Fields = append(errList.Fields, v1.FieldError{
			Field: fmt.Sprintf("emails.%d", i),
			Cause: err,
		})
	}

	if err := errList.Validate(); err != nil {
		span.SetStatus(codes.Error, "failed to add users to ledger")
		return err
	}

	return nil
}
