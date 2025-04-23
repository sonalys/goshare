package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) GetExpense(ctx context.Context, params handlers.GetExpenseParams) (*handlers.Expense, error) {
	_, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	expense, err := a.dependencies.LedgerController.GetExpense(ctx, v1.ConvertID(params.LedgerID), v1.ConvertID(params.ExpenseID))
	if err != nil {
		return nil, err
	}

	return convertExpense(expense), nil
}

func convertExpense(expense *v1.Expense) *handlers.Expense {
	return &handlers.Expense{
		ID:          handlers.NewOptUUID(expense.ID.UUID()),
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate,
		Records:     convertRecords(expense.Records),
	}
}

func convertRecords(records []v1.Record) []handlers.ExpenseRecord {
	result := make([]handlers.ExpenseRecord, 0, len(records))
	for _, record := range records {
		result = append(result, handlers.ExpenseRecord{
			ID:         handlers.NewOptUUID(record.ID.UUID()),
			Type:       handlers.ExpenseRecordType(record.Type.String()),
			FromUserID: record.From.UUID(),
			ToUserID:   record.To.UUID(),
			Amount:     record.Amount,
		})
	}
	return result
}
