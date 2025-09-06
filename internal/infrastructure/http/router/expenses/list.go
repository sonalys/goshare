package expenses

import (
	"context"
	"time"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseList(ctx context.Context, params server.LedgerExpenseListParams) (*server.LedgerExpenseListOK, error) {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	result, err := a.Expenses().List(ctx, usercontroller.ListExpensesRequest{
		ActorID:  identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Limit:    params.Limit.Or(10),
		Cursor:   params.Cursor.Or(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	var cursor server.OptDateTime

	if result.Cursor != nil {
		cursor = server.NewOptDateTime(*result.Cursor)
	}

	return &server.LedgerExpenseListOK{
		Expenses: mapLedgerExpenseToResponseObject(result.Expenses),
		Cursor:   cursor,
	}, nil
}

func mapLedgerExpenseToResponseObject(expenses []v1.LedgerExpenseSummary) []server.ExpenseSummary {
	expensesResponse := make([]server.ExpenseSummary, 0, len(expenses))

	for _, e := range expenses {
		expensesResponse = append(expensesResponse, server.ExpenseSummary{
			ID:          e.ID.UUID(),
			Amount:      e.Amount,
			Name:        e.Name,
			ExpenseDate: e.ExpenseDate,
			CreatedAt:   e.CreatedAt,
			CreatedBy:   e.CreatedBy.UUID(),
			UpdatedAt:   e.UpdatedAt,
			UpdatedBy:   e.UpdatedBy.UUID(),
		})
	}

	return expensesResponse
}
