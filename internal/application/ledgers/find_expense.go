package ledgers

import (
	"context"
	"log/slog"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (c *Controller) FindExpense(ctx context.Context, _ v1.ID, expenseID v1.ID) (*v1.Expense, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.FindExpense")
	defer span.End()

	logFields := []any{
		slog.String("expense_id", expenseID.String()),
	}

	expense, err := c.expenseRepository.Find(ctx, expenseID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get ledger expense", append(logFields, slog.Any("error", err))...)
		return nil, err
	}

	slog.InfoContext(ctx, "ledger expense retrieved")

	return expense, nil
}
