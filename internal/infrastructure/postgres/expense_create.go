package postgres

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *ExpenseRepository) Create(ctx context.Context, ledgerID domain.ID, expense *domain.Expense) error {
	return r.transaction(ctx, func(conn connection) error {
		query := conn.queries()

		createExpenseReq := sqlcgen.CreateExpenseParams{
			ID:          expense.ID,
			LedgerID:    expense.LedgerID,
			Amount:      expense.Amount,
			Name:        expense.Name,
			ExpenseDate: convertTime(expense.ExpenseDate),
			CreatedAt:   convertTime(expense.CreatedAt),
			CreatedBy:   expense.CreatedBy,
			UpdatedAt:   convertTime(expense.UpdatedAt),
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
				CreatedAt:  convertTime(record.CreatedAt),
				CreatedBy:  record.CreatedBy,
				UpdatedAt:  convertTime(record.UpdatedAt),
				UpdatedBy:  record.UpdatedBy,
			}

			if err := query.CreateExpenseRecord(ctx, createRecordReq); err != nil {
				return fmt.Errorf("creating expense record: %w", err)
			}
		}

		return nil
	})
}
