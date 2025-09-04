package expense

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func createExpenseParams(expense *domain.Expense) sqlcgen.CreateExpenseParams {
	return sqlcgen.CreateExpenseParams{
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
}

func createExpenseRecordParams(expenseID, recordID domain.ID, record *domain.Record) sqlcgen.CreateExpenseRecordParams {
	return sqlcgen.CreateExpenseRecordParams{
		ID:         recordID,
		ExpenseID:  expenseID,
		FromUserID: record.From,
		ToUserID:   record.To,
		RecordType: record.Type.String(),
		Amount:     record.Amount,
		CreatedAt:  postgres.ConvertTime(record.CreatedAt),
		CreatedBy:  record.CreatedBy,
		UpdatedAt:  postgres.ConvertTime(record.UpdatedAt),
		UpdatedBy:  record.UpdatedBy,
	}
}

func (r *Repository) Create(ctx context.Context, expense *domain.Expense) error {
	return r.transaction(ctx, func(conn postgres.Connection) error {
		query := conn.Queries()

		if err := query.CreateExpense(ctx, createExpenseParams(expense)); err != nil {
			return fmt.Errorf("creating expense: %w", err)
		}

		for recordID, record := range expense.Records {
			if err := query.CreateExpenseRecord(ctx, createExpenseRecordParams(expense.ID, recordID, record)); err != nil {
				return fmt.Errorf("creating expense record: %w", err)
			}
		}

		return nil
	})
}
