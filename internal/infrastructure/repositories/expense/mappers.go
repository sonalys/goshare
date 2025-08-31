package expense

import (
	"fmt"

	domain "github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
	v1 "github.com/sonalys/goshare/pkg/v1"
)

func toLedgerExpenseSummary(expense *sqlcgen.Expense) *v1.LedgerExpenseSummary {
	return &v1.LedgerExpenseSummary{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate.Time,
		CreatedAt:   expense.CreatedAt.Time,
		CreatedBy:   expense.CreatedBy,
		UpdatedAt:   expense.UpdatedAt.Time,
		UpdatedBy:   expense.UpdatedBy,
	}
}

func toExpense(expense *sqlcgen.Expense, records []sqlcgen.ExpenseRecord) (*domain.Expense, error) {
	result := &domain.Expense{
		ID:          expense.ID,
		LedgerID:    expense.LedgerID,
		Amount:      expense.Amount,
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate.Time,
		CreatedAt:   expense.CreatedAt.Time,
		CreatedBy:   expense.CreatedBy,
		UpdatedAt:   expense.UpdatedAt.Time,
		UpdatedBy:   expense.UpdatedBy,
		Records:     make(map[domain.ID]*domain.Record, len(records)),
	}

	for _, recordModel := range records {
		record, err := toRecord(&recordModel)
		if err != nil {
			return nil, fmt.Errorf("creating record: %w", err)
		}
		result.Records[recordModel.ID] = record
	}

	return result, nil
}

func toRecord(record *sqlcgen.ExpenseRecord) (*domain.Record, error) {
	recordType, err := domain.NewRecordType(record.RecordType)
	if err != nil {
		return nil, fmt.Errorf("invalid record type: %w", err)
	}

	return &domain.Record{
		From:      record.FromUserID,
		To:        record.ToUserID,
		Type:      recordType,
		Amount:    record.Amount,
		CreatedAt: record.CreatedAt.Time,
		CreatedBy: record.CreatedBy,
		UpdatedAt: record.UpdatedAt.Time,
		UpdatedBy: record.UpdatedBy,
	}, nil
}
