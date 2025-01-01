package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

// CreateExpense implements handlers.StrictServerInterface.
func (a *API) CreateExpense(ctx context.Context, request handlers.CreateExpenseRequestObject) (handlers.CreateExpenseResponseObject, error) {
	identity, err := GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	// CreateExpenseRequest is a struct that contains the request parameters.
	req := ledgers.CreateExpenseRequest{
		UserID:       identity.UserID,
		LedgerID:     request.LedgerID,
		CategoryID:   request.Body.CategoryId,
		Amount:       int32(request.Body.Amount),
		Name:         request.Body.Name,
		ExpenseDate:  request.Body.ExpenseDate,
		UserBalances: convertUserBalances(request.Body.UserBalances),
	}

	switch resp, err := a.dependencies.ExpenseCreater.CreateExpense(ctx, req); {
	case err == nil:
		return handlers.CreateExpense200JSONResponse{
			Id: resp.ID,
		}, nil
	default:
		return nil, err
	}
}

func convertUserBalances(userBalances []handlers.ExpenseUserBalance) []v1.ExpenseUserBalance {
	var balances []v1.ExpenseUserBalance
	for _, ub := range userBalances {
		balances = append(balances, v1.ExpenseUserBalance{
			UserID:  ub.UserId,
			Balance: int32(ub.Balance),
		})
	}
	return balances
}
