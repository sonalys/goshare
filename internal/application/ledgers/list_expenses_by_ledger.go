package ledgers

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/pointers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	ListByLedgerParams struct {
		LedgerID v1.ID
		Cursor   *time.Time
		Limit    int32
	}

	ListByLedgerResponse struct {
		Expenses []v1.Expense
		Cursor   *time.Time
	}
)

func (c *Controller) ListExpensesByLedger(ctx context.Context, params ListByLedgerParams) (*ListByLedgerResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.ListExpensesByLedger")
	defer span.End()

	params.Limit = max(params.Limit, 1)

	expenses, err := c.expenseRepository.GetByLedger(ctx, params.LedgerID, pointers.Coalesce(params.Cursor, time.Time{}), params.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list expenses: %w", err)
	}

	var cursor *time.Time

	if len(expenses) == int(params.Limit) {
		cursor = pointers.New(expenses[0].CreatedAt)
		for i := range expenses {
			if cur := expenses[i].CreatedAt; cur.Before(*cursor) {
				*cursor = cur
			}
		}
	}

	return &ListByLedgerResponse{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}
