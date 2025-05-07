package api

import (
	"context"
	"errors"
	"time"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerExpenseList(ctx context.Context, params handlers.LedgerExpenseListParams) (handlers.LedgerExpenseListRes, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	result, err := a.Ledgers.GetExpenses(ctx, controllers.GetExpensesRequest{
		Identity: identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Limit:    params.Limit.Or(10),
		Cursor:   params.Cursor.Or(time.Now()),
	})
	switch {
	case err == nil:
		var cursor handlers.OptDateTime

		if result.Cursor != nil {
			cursor = handlers.NewOptDateTime(*result.Cursor)
		}

		return &handlers.LedgerExpenseListOK{
			Expenses: mapLedgerExpenseToResponseObject(result.Expenses),
			Cursor:   cursor,
		}, nil
	case errors.Is(err, domain.ErrUserNotAMember):
		return newRespUnauthorized(ctx), nil
	default:
		return nil, err
	}
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
