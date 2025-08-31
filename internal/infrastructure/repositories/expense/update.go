package expense

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func updateExpenseParams(expense *domain.Expense) sqlcgen.UpdateExpenseParams {
	return sqlcgen.UpdateExpenseParams{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Name:        expense.Name,
		ExpenseDate: postgres.ConvertTime(expense.ExpenseDate),
		UpdatedAt:   postgres.ConvertTime(expense.UpdatedAt),
		UpdatedBy:   expense.UpdatedBy,
	}
}

func (r *Repository) Update(ctx context.Context, expense *domain.Expense) error {
	return r.transaction(ctx, func(conn postgres.Connection) error {
		query := conn.Queries()

		if err := query.DeleteExpenseRecordsNotIn(ctx, slices.Collect(maps.Keys(expense.Records))); err != nil {
			return fmt.Errorf("deleting records: %w", err)
		}

		if ids, err := query.UpdateExpense(ctx, updateExpenseParams(expense)); err != nil {
			return fmt.Errorf("updating expense: %w", err)
		} else if len(ids) == 0 {
			return domain.ErrExpenseNotFound
		}

		for id, record := range expense.Records {
			if err := query.CreateExpenseRecord(ctx, createExpenseRecordParams(expense.ID, id, record)); err != nil {
				return fmt.Errorf("creating record: %w", err)
			}
		}

		return nil
	})
}
