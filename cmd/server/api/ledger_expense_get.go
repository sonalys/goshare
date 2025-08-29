package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerExpenseGet(ctx context.Context, params handlers.LedgerExpenseGetParams) (*handlers.Expense, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	expense, err := a.UserController.FindExpense(ctx, usercontroller.FindExpenseRequest{
		Actor:     identity.UserID,
		LedgerID:  domain.ConvertID(params.LedgerID),
		ExpenseID: domain.ConvertID(params.ExpenseID),
	})
	if err != nil {
		return nil, err
	}

	return convertExpense(expense), nil
}

func convertExpense(expense *domain.Expense) *handlers.Expense {
	return &handlers.Expense{
		ID:          handlers.NewOptUUID(expense.ID.UUID()),
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate,
		Records:     convertRecords(expense.Records),
	}
}

func convertRecords(records map[domain.ID]*domain.Record) []handlers.ExpenseRecord {
	result := make([]handlers.ExpenseRecord, 0, len(records))
	for id, record := range records {
		result = append(result, handlers.ExpenseRecord{
			ID:         handlers.NewOptUUID(id.UUID()),
			Type:       handlers.ExpenseRecordType(record.Type.String()),
			FromUserID: record.From.UUID(),
			ToUserID:   record.To.UUID(),
			Amount:     record.Amount,
		})
	}
	return result
}
