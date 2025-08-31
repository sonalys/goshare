package usercontroller

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/slog"
	v1 "github.com/sonalys/goshare/pkg/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DeleteExpenseRecordRequest struct {
	ActorID   domain.ID
	LedgerID  domain.ID
	ExpenseID domain.ID
	RecordID  domain.ID
}

func (c *recordsController) Delete(ctx context.Context, req DeleteExpenseRecordRequest) error {
	ctx, span := c.tracer.Start(ctx, "delete",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
			attribute.Stringer("ledger_id", req.LedgerID),
			attribute.Stringer("expense_id", req.ExpenseID),
			attribute.Stringer("record_id", req.RecordID),
		),
	)
	defer span.End()

	err := c.db.Transaction(ctx, func(db ports.LocalRepositories) error {
		ledger, err := db.Ledger().Get(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("getting ledger: %w", err)
		}

		if !ledger.CanManageExpenses(req.ActorID) {
			return fmt.Errorf("authorizing expenses management: %w", v1.ErrForbidden)
		}

		expense, err := db.Expense().Get(ctx, req.ExpenseID)
		if err != nil {
			return fmt.Errorf("getting expense: %w", err)
		}

		if err = expense.DeleteRecord(ledger, req.RecordID); err != nil {
			return fmt.Errorf("deleting record: %w", err)
		}

		if err = db.Ledger().Update(ctx, ledger); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		if err = db.Expense().Update(ctx, expense); err != nil {
			return fmt.Errorf("updating expense: %w", err)
		}

		return nil
	})
	if err != nil {
		return slog.ErrorReturn(ctx, "committing transaction", err)
	}

	slog.Info(ctx, "ledger expense record deleted")

	return nil
}
