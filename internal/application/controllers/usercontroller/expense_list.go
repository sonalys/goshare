package usercontroller

import (
	"context"
	"time"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	ListExpensesRequest struct {
		Actor    domain.ID
		LedgerID domain.ID
		Cursor   time.Time
		Limit    int32
	}

	ListExpensesResponse struct {
		Expenses []v1.LedgerExpenseSummary
		Cursor   *time.Time
	}
)

func (c *expenseController) List(ctx context.Context, req ListExpensesRequest) (*ListExpensesResponse, error) {
	ctx, span := c.tracer.Start(ctx, "list",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	req.Limit = max(1, req.Limit)

	expenses, err := c.db.Expense().GetByLedger(ctx, req.LedgerID, req.Cursor, req.Limit+1)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "listing expenses", err)
	}

	slog.Info(ctx, "ledger expenses retrieved")

	if len(expenses) == 0 {
		return &ListExpensesResponse{}, nil
	}

	var cursor *time.Time
	if len(expenses) == int(req.Limit)+1 {
		expenses = expenses[:len(expenses)-1]
		cursor = &expenses[len(expenses)-1].CreatedAt
	}

	return &ListExpensesResponse{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}
