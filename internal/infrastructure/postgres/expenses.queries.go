package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *ExpensesRepository) createExpense(ctx context.Context, tx pgx.Tx, expense *v1.Expense) error {
	createExpenseReq := sqlc.CreateExpenseParams{
		ID:          convertUUID(expense.ID),
		Amount:      expense.Amount,
		CategoryID:  convertUUIDPtr(expense.CategoryID),
		LedgerID:    convertUUID(expense.LedgerID),
		Name:        expense.Name,
		ExpenseDate: convertTime(expense.ExpenseDate),
		CreatedAt:   convertTime(expense.CreatedAt),
		CreatedBy:   convertUUID(expense.CreatedBy),
		UpdatedAt:   convertTime(expense.UpdatedAt),
		UpdatedBy:   convertUUID(expense.UpdatedBy),
	}

	queries := r.client.queries().WithTx(tx)

	if err := queries.CreateExpense(ctx, createExpenseReq); err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	for i, balance := range expense.UserBalances {
		ledgerRecord := sqlc.AppendLedgerRecordParams{
			ID:          convertUUID(v1.NewID()),
			LedgerID:    convertUUID(expense.LedgerID),
			ExpenseID:   convertUUID(expense.ID),
			UserID:      convertUUID(balance.UserID),
			Amount:      balance.Balance,
			Description: "Initial balance for new expense",
			CreatedAt:   convertTime(expense.CreatedAt),
			CreatedBy:   convertUUID(expense.CreatedBy),
		}

		if err := queries.AppendLedgerRecord(ctx, ledgerRecord); err != nil {
			return fmt.Errorf("failed to append ledger record %d: %w", i, err)
		}
	}

	return nil
}
