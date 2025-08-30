package postgres

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
)

func (r *ExpenseRepository) Get(ctx context.Context, id domain.ID) (*domain.Expense, error) {
	expense, err := r.client.queries().GetExpenseById(ctx, id)
	if err != nil {
		return nil, mapLedgerError(err)
	}

	records, err := r.client.queries().GetExpenseRecords(ctx, expense.ID)
	if err != nil {
		return nil, fmt.Errorf("getting expense records: %w", err)
	}

	return mappers.NewExpense(&expense, records)
}
