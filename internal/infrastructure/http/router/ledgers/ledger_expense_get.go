package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseGet(ctx context.Context, params server.LedgerExpenseGetParams) (*server.Expense, error) {
	identity, err := middlewares.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	expense, err := a.UserController.Expenses().Get(ctx, usercontroller.GetExpenseRequest{
		ActorID:   identity.UserID,
		LedgerID:  domain.ConvertID(params.LedgerID),
		ExpenseID: domain.ConvertID(params.ExpenseID),
	})
	if err != nil {
		return nil, err
	}

	return convertExpense(expense), nil
}

func convertExpense(expense *domain.Expense) *server.Expense {
	return &server.Expense{
		ID:          server.NewOptUUID(expense.ID.UUID()),
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate,
		Records:     convertRecords(expense.Records),
	}
}

func convertRecords(records map[domain.ID]*domain.Record) []server.ExpenseRecord {
	result := make([]server.ExpenseRecord, 0, len(records))
	for id, record := range records {
		result = append(result, server.ExpenseRecord{
			ID:         server.NewOptUUID(id.UUID()),
			Type:       server.ExpenseRecordType(record.Type.String()),
			FromUserID: record.From.UUID(),
			ToUserID:   record.To.UUID(),
			Amount:     record.Amount,
		})
	}
	return result
}
