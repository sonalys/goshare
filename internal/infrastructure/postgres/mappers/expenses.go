package mappers

import (
	"fmt"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	domain "github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func NewLedgerExpenseSummary(expense *sqlc.Expense) *v1.LedgerExpenseSummary {
	return &v1.LedgerExpenseSummary{
		ID:          newUUID(expense.ID),
		Amount:      expense.Amount,
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate.Time,
		CreatedAt:   expense.CreatedAt.Time,
		CreatedBy:   newUUID(expense.CreatedBy),
		UpdatedAt:   expense.UpdatedAt.Time,
		UpdatedBy:   newUUID(expense.UpdatedBy),
	}
}

func NewExpense(expense *sqlc.Expense, records []sqlc.ExpenseRecord) (*domain.Expense, error) {
	result := &domain.Expense{
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
		record, err := NewRecord(&record)
		if err != nil {
			return nil, fmt.Errorf("creating record: %w", err)
		}
		result.Records = append(result.Records, *record)
	}

	return result, nil
}

func NewRecord(record *sqlc.ExpenseRecord) (*domain.Record, error) {
	recordType, err := domain.NewRecordType(record.RecordType)
	if err != nil {
		return nil, fmt.Errorf("invalid record type: %w", err)
	}
	return &domain.Record{
		ID:        newUUID(record.ID),
		From:      newUUID(record.FromUserID),
		To:        newUUID(record.ToUserID),
		Type:      recordType,
		Amount:    record.Amount,
		CreatedAt: record.CreatedAt.Time,
		CreatedBy: newUUID(record.CreatedBy),
		UpdatedAt: record.UpdatedAt.Time,
		UpdatedBy: newUUID(record.UpdatedBy),
	}, nil
}
