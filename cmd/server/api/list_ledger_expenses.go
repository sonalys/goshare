package api

import (
	"context"
	"time"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) ListLedgerExpenses(ctx context.Context, params handlers.ListLedgerExpensesParams) (r *handlers.ListLedgerExpensesOK, _ error) {
	var cursor *time.Time
	if t, ok := params.Cursor.Get(); ok {
		cursor = &t
	}
	apiParams := ledgers.ListByLedgerParams{
		LedgerID: v1.ConvertID(params.LedgerID),
		Cursor:   cursor,
		Limit:    params.Limit.Value,
	}

	switch resp, err := a.dependencies.LedgerController.ListExpensesByLedger(ctx, apiParams); {
	case err == nil:
		var cursor handlers.OptDateTime
		if resp.Cursor != nil {
			cursor = handlers.NewOptDateTime(*resp.Cursor)
		}
		return &handlers.ListLedgerExpensesOK{
			Expenses: convertExpenses(resp.Expenses),
			Cursor:   cursor,
		}, nil
	default:
		return nil, err
	}
}

func convertExpenseUserBalances(from []v1.ExpenseUserBalance) []handlers.ExpenseUserBalance {
	to := make([]handlers.ExpenseUserBalance, 0, len(from))

	for i := range from {
		to = append(to, handlers.ExpenseUserBalance{
			Balance: from[i].Balance,
			UserID:  from[i].UserID.UUID(),
		})
	}

	return to
}

func convertExpenses(from []v1.Expense) []handlers.LedgerExpense {
	to := make([]handlers.LedgerExpense, 0, len(from))

	for i := range from {
		cur := &from[i]

		var categoryID handlers.OptUUID
		if cur.CategoryID != nil {
			categoryID = handlers.NewOptUUID(cur.CategoryID.UUID())
		}

		to = append(to, handlers.LedgerExpense{
			ID:           cur.ID.UUID(),
			CategoryID:   categoryID,
			ExpenseDate:  cur.ExpenseDate,
			Name:         cur.Name,
			UserBalances: convertExpenseUserBalances(cur.UserBalances),
			Amount:       cur.Amount,
			CreatedAt:    cur.CreatedAt,
			CreatedBy:    cur.CreatedBy.UUID(),
			UpdatedAt:    cur.UpdatedAt,
			UpdatedBy:    cur.UpdatedBy.UUID(),
		})
	}

	return to
}
