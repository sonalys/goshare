package usercontroller

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/slog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	CreateLedgerRequest struct {
		ActorID domain.ID
		Name    string
	}

	CreateLedgerResponse struct {
		ID domain.ID
	}
)

func (c *ledgerController) Create(ctx context.Context, req CreateLedgerRequest) (resp *CreateLedgerResponse, err error) {
	ctx, span := c.tracer.Start(ctx, "create",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
		),
	)
	defer span.End()

	slog.Debug(ctx, "creating ledger", slog.With("req", req))

	err = c.db.Transaction(ctx, func(db ports.LocalRepositories) error {
		user, err := db.User().Get(ctx, req.ActorID)
		if err != nil {
			return fmt.Errorf("finding user: %w", err)
		}

		ledger, err := user.CreateLedger(req.Name)
		if err != nil {
			return fmt.Errorf("creating ledger: %w", err)
		}

		if err := db.Ledger().Create(ctx, ledger); err != nil {
			return fmt.Errorf("saving ledger: %w", err)
		}

		resp = &CreateLedgerResponse{
			ID: ledger.ID,
		}

		if err := db.User().Create(ctx, user); err != nil {
			return fmt.Errorf("saving user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "committing transaction", err)
	}

	slog.Info(ctx, "ledger created", slog.WithStringer("ledger_id", resp.ID))

	return
}
