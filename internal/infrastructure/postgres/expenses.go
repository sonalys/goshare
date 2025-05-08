package postgres

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type ExpenseRepository struct {
	client connection
}

func NewExpenseRepository(client connection) *ExpenseRepository {
	return &ExpenseRepository{
		client: client,
	}
}

func (r *ExpenseRepository) Create(ctx context.Context, ledgerID domain.ID, expense *domain.Expense) error {
	return r.client.transaction(ctx, func(query *sqlc.Queries) error {
		createExpenseReq := sqlc.CreateExpenseParams{
			ID:          convertID(expense.ID),
			LedgerID:    convertID(expense.LedgerID),
			Amount:      expense.Amount,
			Name:        expense.Name,
			ExpenseDate: convertTime(expense.ExpenseDate),
			CreatedAt:   convertTime(expense.CreatedAt),
			CreatedBy:   convertID(expense.CreatedBy),
			UpdatedAt:   convertTime(expense.UpdatedAt),
			UpdatedBy:   convertID(expense.UpdatedBy),
		}

		if err := query.CreateExpense(ctx, createExpenseReq); err != nil {
			return fmt.Errorf("creating expense: %w", err)
		}

		for _, record := range expense.Records {
			createRecordReq := sqlc.CreateExpenseRecordParams{
				ID:         convertID(record.ID),
				ExpenseID:  convertID(expense.ID),
				FromUserID: convertID(record.From),
				ToUserID:   convertID(record.To),
				RecordType: record.Type.String(),
				Amount:     record.Amount,
				CreatedAt:  convertTime(record.CreatedAt),
				CreatedBy:  convertID(record.CreatedBy),
				UpdatedAt:  convertTime(record.UpdatedAt),
				UpdatedBy:  convertID(record.UpdatedBy),
			}

			if err := query.CreateExpenseRecord(ctx, createRecordReq); err != nil {
				return fmt.Errorf("creating expense record: %w", err)
			}
		}

		return nil
	})
}

func (r *ExpenseRepository) Find(ctx context.Context, id domain.ID) (*domain.Expense, error) {
	expense, err := r.client.queries().FindExpenseById(ctx, convertID(id))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	records, err := r.client.queries().GetExpenseRecords(ctx, expense.ID)
	if err != nil {
		return nil, fmt.Errorf("getting expense records: %w", err)
	}

	return mappers.NewExpense(&expense, records)
}

func (r *ExpenseRepository) GetByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error) {
	expenses, err := r.client.queries().GetLedgerExpenses(ctx, sqlc.GetLedgerExpensesParams{
		LedgerID:  convertID(ledgerID),
		Limit:     limit,
		CreatedAt: convertTime(cursor),
	})
	if err != nil {
		return nil, fmt.Errorf("getting ledger expenses: %w", err)
	}

	result := make([]v1.LedgerExpenseSummary, 0, len(expenses))
	for _, expense := range expenses {
		result = append(result, *mappers.NewLedgerExpenseSummary(&expense))
	}

	return result, nil
}
