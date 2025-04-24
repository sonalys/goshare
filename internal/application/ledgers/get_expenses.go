package ledgers

import (
	"context"
	"log/slog"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type GetExpensesParams struct {
	LedgerID v1.ID
	Cursor   time.Time
	Limit    int32
}

type GetExpensesResult struct {
	Expenses []v1.LedgerExpenseSummary
	Cursor   *time.Time
}

func (c *Controller) GetExpenses(ctx context.Context, params GetExpensesParams) (*GetExpensesResult, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetExpenses")
	defer span.End()

	params.Limit = max(1, params.Limit)

	logFields := []any{
		slog.String("ledger_id", params.LedgerID.String()),
	}

	expenses, err := c.expenseRepository.GetByLedger(ctx, params.LedgerID, params.Cursor, params.Limit+1)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get ledger expenses", append(logFields, slog.Any("error", err))...)
		return nil, err
	}

	if len(expenses) == 0 {
		return nil, nil
	}

	slog.InfoContext(ctx, "ledger expenses retrieved", logFields...)

	var cursor *time.Time
	if len(expenses) == int(params.Limit)+1 {
		expenses = expenses[:len(expenses)-1]
		cursor = &expenses[len(expenses)-1].CreatedAt
	}

	return &GetExpensesResult{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}
