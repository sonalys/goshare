package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type ExpenseRepository struct {
	client *Client
}

func NewExpenseRepository(client *Client) *ExpenseRepository {
	return &ExpenseRepository{
		client: client,
	}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense *v1.Expense) error {
	r.client.transaction(ctx, func(tx pgx.Tx) error {
		query := r.client.queries().WithTx(tx)

		createExpenseReq := sqlc.CreateExpenseParams{
			ID:          convertUUID(expense.ID),
			LedgerID:    convertUUID(expense.LedgerID),
			Amount:      expense.Amount,
			Name:        expense.Name,
			ExpenseDate: convertTime(expense.ExpenseDate),
			CreatedAt:   convertTime(expense.CreatedAt),
			CreatedBy:   convertUUID(expense.CreatedBy),
			UpdatedAt:   convertTime(expense.UpdatedAt),
			UpdatedBy:   convertUUID(expense.UpdatedBy),
		}

		if err := query.CreateExpense(ctx, createExpenseReq); err != nil {
			return fmt.Errorf("creating expense: %w", err)
		}

		for _, record := range expense.Records {
			createRecordReq := sqlc.CreateExpenseRecordParams{
				ID:         convertUUID(record.ID),
				ExpenseID:  convertUUID(expense.ID),
				FromUserID: convertUUID(record.From),
				ToUserID:   convertUUID(record.To),
				RecordType: record.Type.String(),
				Amount:     record.Amount,
				CreatedAt:  convertTime(record.CreatedAt),
				CreatedBy:  convertUUID(record.CreatedBy),
				UpdatedAt:  convertTime(record.UpdatedAt),
				UpdatedBy:  convertUUID(record.UpdatedBy),
			}

			if err := query.CreateExpenseRecord(ctx, createRecordReq); err != nil {
				return fmt.Errorf("creating expense record: %w", err)
			}
		}

		return nil
	})
	return nil
}

func (r *ExpenseRepository) Find(ctx context.Context, id v1.ID) (*v1.Expense, error) {
	expense, err := r.client.queries().FindExpenseById(ctx, convertUUID(id))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	records, err := r.client.queries().GetExpenseRecords(ctx, expense.ID)
	if err != nil {
		return nil, fmt.Errorf("getting expense records: %w", err)
	}

	return NewExpense(&expense, records), nil
}

func NewExpense(expense *sqlc.Expense, records []sqlc.ExpenseRecord) *v1.Expense {
	result := &v1.Expense{
		ID:          newUUID(expense.ID),
		LedgerID:    newUUID(expense.LedgerID),
		Amount:      expense.Amount,
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate.Time,
		CreatedAt:   expense.CreatedAt.Time,
		CreatedBy:   newUUID(expense.CreatedBy),
		UpdatedAt:   expense.UpdatedAt.Time,
		UpdatedBy:   newUUID(expense.UpdatedBy),
	}

	for _, record := range records {
		result.Records = append(result.Records, *NewRecord(&record))
	}

	return result
}

func NewRecord(record *sqlc.ExpenseRecord) *v1.Record {
	return &v1.Record{
		ID:        newUUID(record.ID),
		From:      newUUID(record.FromUserID),
		To:        newUUID(record.ToUserID),
		Type:      v1.NewRecordType(record.RecordType),
		Amount:    record.Amount,
		CreatedAt: record.CreatedAt.Time,
		CreatedBy: newUUID(record.CreatedBy),
		UpdatedAt: record.UpdatedAt.Time,
		UpdatedBy: newUUID(record.UpdatedBy),
	}
}
