package expenses

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseGet(ctx context.Context, params server.LedgerExpenseGetParams) (*server.Expense, error) {
	identity, err := a.securityHandler.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	expense, err := a.controller.Expenses().Get(ctx, usercontroller.GetExpenseRequest{
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
		Records:     mappers.RecordToExpenseRecord(expense.Records),
	}
}
