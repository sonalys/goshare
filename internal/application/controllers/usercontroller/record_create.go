package usercontroller

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CreateExpenseRecordRequest struct {
	ActorID        domain.ID
	LedgerID       domain.ID
	ExpenseID      domain.ID
	PendingRecords []domain.PendingRecord
}

func (c *recordsController) Create(ctx context.Context, req CreateExpenseRecordRequest) (resp *domain.Expense, err error) {
	ctx, span := c.tracer.Start(ctx, "create",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
			attribute.Stringer("ledger_id", req.LedgerID),
			attribute.Stringer("expense_id", req.ExpenseID),
		),
	)
	defer span.End()

	slog.Debug(ctx, "creating expense record", slog.With("req", req))

	err = c.db.Transaction(ctx, func(db application.Repositories) error {
		ledger, err := db.Ledger().Get(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("fetching ledger: %w", err)
		}

		if !ledger.CanManageExpenses(req.ActorID) {
			return fmt.Errorf("authorizing expenses management: %w", application.ErrUnauthorized)
		}

		expense, err := db.Expense().Get(ctx, req.ExpenseID)
		if err != nil {
			return fmt.Errorf("fetching expense: %w", err)
		}

		if err := expense.CreateRecords(req.ActorID, ledger, req.PendingRecords...); err != nil {
			return fmt.Errorf("appending new records: %w", err)
		}

		if err := db.Ledger().Update(ctx, ledger); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		if err := db.Expense().Update(ctx, expense); err != nil {
			return fmt.Errorf("updating expense: %w", err)
		}

		resp = expense
		return nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "commiting transaction", err)
	}

	slog.Info(ctx, "expense records created")

	return
}
