package usercontroller

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetExpenseRequest struct {
	ActorID   domain.ID
	LedgerID  domain.ID
	ExpenseID domain.ID
}

func (c *expenseController) Get(ctx context.Context, req GetExpenseRequest) (*domain.Expense, error) {
	ctx, span := c.tracer.Start(ctx, "get",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
			attribute.Stringer("ledger_id", req.LedgerID),
			attribute.Stringer("expense_id", req.ExpenseID),
		),
	)
	defer span.End()

	ledger, err := c.db.Ledger().Get(ctx, req.LedgerID)
	if err != nil {
		return nil, fmt.Errorf("finding ledger: %w", err)
	}

	if !ledger.CanView(req.ActorID) {
		return nil, fmt.Errorf("authorizing user ledger view: %w", v1.ErrForbidden)
	}

	expense, err := c.db.Expense().Get(ctx, req.ExpenseID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "getting expense", err)
	}

	slog.Info(ctx, "ledger expense retrieved")

	return expense, nil
}
