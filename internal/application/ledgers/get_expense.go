package ledgers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (c *Controller) GetExpense(ctx context.Context, _ v1.ID, expenseID v1.ID) (*v1.Expense, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetExpense")
	defer span.End()

	expense, err := c.expenseRepository.Find(ctx, expenseID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get ledger expense", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get ledger expense: %w", err)
	}

	slog.InfoContext(ctx, "ledger expense retrieved")

	return expense, nil
}
