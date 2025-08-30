package usercontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	CreateExpenseRequest struct {
		Actor          domain.ID
		LedgerID       domain.ID
		Name           string
		ExpenseDate    time.Time
		PendingRecords []domain.PendingRecord
	}

	CreateExpenseResponse struct {
		ID domain.ID
	}
)

func (c *expenseController) Create(ctx context.Context, req CreateExpenseRequest) (resp *CreateExpenseResponse, err error) {
	ctx, span := c.tracer.Start(ctx, "create",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	slog.Debug(ctx, "creating expense", slog.With("req", req))

	err = c.db.Transaction(ctx, func(db application.Repositories) error {
		ledger, err := db.Ledger().Get(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("finding ledger: %w", err)
		}

		expense, err := ledger.CreateExpense(domain.CreateExpenseRequest{
			Creator:        req.Actor,
			Name:           req.Name,
			ExpenseDate:    req.ExpenseDate,
			PendingRecords: req.PendingRecords,
		})
		if err != nil {
			return fmt.Errorf("creating expense: %w", err)
		}

		if err = db.Expense().Create(ctx, req.LedgerID, expense); err != nil {
			return fmt.Errorf("saving expense: %w", err)
		}

		if err = db.Ledger().Update(ctx, ledger); err != nil {
			return fmt.Errorf("saving ledger: %w", err)
		}

		resp = &CreateExpenseResponse{
			ID: expense.ID,
		}

		return nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "creating expense", err)
	}
	slog.Info(ctx, "expense created")
	return
}
