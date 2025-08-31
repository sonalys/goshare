package expense

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
)

func (r *Repository) Get(ctx context.Context, id domain.ID) (*domain.Expense, error) {
	expense, err := r.conn.Queries().GetExpenseById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting expense: %w", postgres.DefaultErrorMapping(err))
	}

	records, err := r.conn.Queries().GetExpenseRecords(ctx, expense.ID)
	if err != nil {
		return nil, fmt.Errorf("getting expense records: %w", postgres.DefaultErrorMapping(err))
	}

	return mappers.NewExpense(&expense, records)
}
