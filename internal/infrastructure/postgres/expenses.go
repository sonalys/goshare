package postgres

import (
	"context"
	"fmt"
	"maps"
	"slices"
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
	return r.client.transaction(ctx, func(conn connection) error {
		query := conn.queries()

		createExpenseReq := sqlc.CreateExpenseParams{
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
			createRecordReq := sqlc.CreateExpenseRecordParams{
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

func (r *ExpenseRepository) ListByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error) {
	expenses, err := r.client.queries().GetLedgerExpenses(ctx, sqlc.GetLedgerExpensesParams{
		LedgerID:  ledgerID,
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

func (r *ExpenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	return r.client.transaction(ctx, func(conn connection) error {
		query := conn.queries()

		if err := query.DeleteExpenseRecordsNotIn(ctx, slices.Collect(maps.Keys(expense.Records))); err != nil {
			return err
		}

		err := query.UpdateExpense(ctx, sqlc.UpdateExpenseParams{
			ID:          expense.ID,
			Amount:      expense.Amount,
			Name:        expense.Name,
			ExpenseDate: convertTime(expense.ExpenseDate),
			UpdatedAt:   convertTime(expense.UpdatedAt),
			UpdatedBy:   expense.UpdatedBy,
		})
		if err != nil {
			return err
		}

		for id, record := range expense.Records {
			if err := query.CreateExpenseRecord(ctx, sqlc.CreateExpenseRecordParams{
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
				return err
			}
		}
		return nil
	})
}
