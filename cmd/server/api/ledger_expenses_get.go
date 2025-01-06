package api

import (
	"context"

	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	"github.com/sonalys/goshare/internal/pkg/pointers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) ListLedgerExpenses(ctx context.Context, request handlers.ListLedgerExpensesRequestObject) (handlers.ListLedgerExpensesResponseObject, error) {
	params := ledgers.ListByLedgerParams{
		LedgerID: v1.ConvertID(request.LedgerID),
		Cursor:   request.Params.Cursor,
		Limit:    pointers.Coalesce(request.Params.Limit, 50),
	}

	switch resp, err := a.dependencies.ExpensesLister.ListExpensesByLedger(ctx, params); {
	case err == nil:
		return handlers.ListLedgerExpenses200JSONResponse{
			Expenses: convertExpenses(resp.Expenses),
			Cursor:   resp.Cursor,
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
			UserId:  from[i].UserID.UUID(),
		})
	}

	return to
}

func convertExpenses(from []v1.Expense) []handlers.LedgerExpense {
	to := make([]handlers.LedgerExpense, 0, len(from))

	for i := range from {
		cur := &from[i]

		to = append(to, handlers.LedgerExpense{
			Id:           cur.ID.UUID(),
			CategoryId:   pointers.Convert(cur.CategoryID, func(from v1.ID) types.UUID { return from.UUID() }),
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
