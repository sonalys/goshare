package postgres

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *ExpenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	return r.transaction(ctx, func(conn connection) error {
		query := conn.queries()

		if err := query.DeleteExpenseRecordsNotIn(ctx, slices.Collect(maps.Keys(expense.Records))); err != nil {
			return fmt.Errorf("deleting records: %w", err)
		}

		if err := query.UpdateExpense(ctx, sqlcgen.UpdateExpenseParams{
			ID:          expense.ID,
			Amount:      expense.Amount,
			Name:        expense.Name,
			ExpenseDate: convertTime(expense.ExpenseDate),
			UpdatedAt:   convertTime(expense.UpdatedAt),
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
				CreatedAt:  convertTime(record.CreatedAt),
				CreatedBy:  record.CreatedBy,
				UpdatedAt:  convertTime(record.UpdatedAt),
				UpdatedBy:  record.UpdatedBy,
			}); err != nil {
				return fmt.Errorf("creating record: %w", err)
			}
		}
		return nil
	})
}
