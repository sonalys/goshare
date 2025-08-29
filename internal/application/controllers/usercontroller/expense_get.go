package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetExpenseRequest struct {
	Actor     domain.ID
	LedgerID  domain.ID
	ExpenseID domain.ID
}

func (c *expenseController) Get(ctx context.Context, req GetExpenseRequest) (*domain.Expense, error) {
	ctx, span := c.tracer.Start(ctx, "get",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
			attribute.Stringer("expense_id", req.ExpenseID),
		),
	)
	defer span.End()

	expense, err := c.db.Expense().Find(ctx, req.ExpenseID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "getting expense", err)
	}

	slog.Info(ctx, "ledger expense retrieved")

	return expense, nil
}
