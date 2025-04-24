package api

import (
	"context"
	"time"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) ListExpenses(ctx context.Context, params handlers.ListExpensesParams) (*handlers.ListExpensesOK, error) {
	result, err := a.dependencies.LedgerController.GetExpenses(ctx, ledgers.GetExpensesParams{
		LedgerID: v1.ConvertID(params.LedgerID),
		Limit:    params.Limit.Or(10),
		Cursor:   params.Cursor.Or(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	var cursor handlers.OptDateTime

	if result.Cursor != nil {
		cursor = handlers.NewOptDateTime(*result.Cursor)
	}

	return &handlers.ListExpensesOK{
		Expenses: mapLedgerExpenseToResponseObject(result.Expenses),
		Cursor:   cursor,
	}, nil
}

func mapLedgerExpenseToResponseObject(expenses []v1.LedgerExpenseSummary) []handlers.ExpenseSummary {
	expensesResponse := make([]handlers.ExpenseSummary, 0, len(expenses))

	for _, e := range expenses {
		expensesResponse = append(expensesResponse, handlers.ExpenseSummary{
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
