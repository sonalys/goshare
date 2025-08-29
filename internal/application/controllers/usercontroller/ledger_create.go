package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	CreateLedgerRequest struct {
		Actor domain.ID
		Name  string
	}

	CreateLedgerResponse struct {
		ID domain.ID
	}
)

func (c *ledgerController) Create(ctx context.Context, req CreateLedgerRequest) (resp *CreateLedgerResponse, err error) {
	ctx, span := c.tracer.Start(ctx, "create",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
		),
	)
	defer span.End()

	slog.Debug(ctx, "creating ledger", slog.With("req", req))

	err = c.db.Transaction(ctx, func(db application.Repositories) error {
		user, err := db.User().Find(ctx, req.Actor)
		if err != nil {
			return err
		}

		ledger, err := user.CreateLedger(req.Name)
		if err != nil {
			return err
		}

		if err := db.Ledger().Create(ctx, ledger); err != nil {
			return err
		}

		resp = &CreateLedgerResponse{
			ID: ledger.ID,
		}

		return db.User().Save(ctx, user)
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "creating ledger", err)
	}

	slog.Info(ctx, "ledger created", slog.WithStringer("ledger_id", resp.ID))

	return
}
