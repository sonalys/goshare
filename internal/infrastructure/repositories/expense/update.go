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

func (r *Repository) Update(ctx context.Context, expense *domain.Expense) error {
	return r.transaction(ctx, func(conn postgres.Connection) error {
		query := conn.Queries()

		if err := query.DeleteExpenseRecordsNotIn(ctx, slices.Collect(maps.Keys(expense.Records))); err != nil {
			return fmt.Errorf("deleting records: %w", err)
		}

		if err := query.UpdateExpense(ctx, sqlcgen.UpdateExpenseParams{
			ID:          expense.ID,
			Amount:      expense.Amount,
			Name:        expense.Name,
			ExpenseDate: postgres.ConvertTime(expense.ExpenseDate),
			UpdatedAt:   postgres.ConvertTime(expense.UpdatedAt),
			UpdatedBy:   expense.UpdatedBy,
		}); err != nil {
			return fmt.Errorf("updating expense: %w", err)
		}

		for id, record := range expense.Records {
			if err := query.CreateExpenseRecord(ctx, sqlcgen.CreateExpenseRecordParams{
				ID:         id,
				ExpenseID:  expense.ID,
				RecordType: record.Type.String(),
				Amount:     record.Amount,
				FromUserID: record.From,
				ToUserID:   record.To,
				CreatedAt:  postgres.ConvertTime(record.CreatedAt),
				CreatedBy:  record.CreatedBy,
				UpdatedAt:  postgres.ConvertTime(record.UpdatedAt),
				UpdatedBy:  record.UpdatedBy,
			}); err != nil {
				return fmt.Errorf("creating record: %w", err)
			}
		}

		return nil
	})
}
