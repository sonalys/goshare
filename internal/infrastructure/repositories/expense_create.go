package repositories

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *ExpenseRepository) Create(ctx context.Context, ledgerID domain.ID, expense *domain.Expense) error {
	return r.transaction(ctx, func(conn postgres.Connection) error {
		query := conn.Queries()

		createExpenseReq := sqlcgen.CreateExpenseParams{
			ID:          expense.ID,
			LedgerID:    expense.LedgerID,
			Amount:      expense.Amount,
			Name:        expense.Name,
			ExpenseDate: postgres.ConvertTime(expense.ExpenseDate),
			CreatedAt:   postgres.ConvertTime(expense.CreatedAt),
			CreatedBy:   expense.CreatedBy,
			UpdatedAt:   postgres.ConvertTime(expense.UpdatedAt),
			UpdatedBy:   expense.UpdatedBy,
		}

		if err := query.CreateExpense(ctx, createExpenseReq); err != nil {
			return fmt.Errorf("creating expense: %w", err)
		}

		for id, record := range expense.Records {
			createRecordReq := sqlcgen.CreateExpenseRecordParams{
				ID:         id,
				ExpenseID:  expense.ID,
				FromUserID: record.From,
				ToUserID:   record.To,
				RecordType: record.Type.String(),
				Amount:     record.Amount,
				CreatedAt:  postgres.ConvertTime(record.CreatedAt),
				CreatedBy:  record.CreatedBy,
				UpdatedAt:  postgres.ConvertTime(record.UpdatedAt),
				UpdatedBy:  record.UpdatedBy,
			}

			if err := query.CreateExpenseRecord(ctx, createRecordReq); err != nil {
				return fmt.Errorf("creating expense record: %w", err)
			}
		}

		return nil
	})
}
