package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) CreateExpense(ctx context.Context, request handlers.CreateExpenseRequestObject) (handlers.CreateExpenseResponseObject, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	req := ledgers.CreateExpenseRequest{
		UserID:       identity.UserID,
		LedgerID:     request.LedgerID,
		CategoryID:   request.Body.CategoryId,
		Amount:       request.Body.Amount,
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
		if errList := new(v1.FieldErrorList); errors.As(err, errList) {
			return handlers.CreateExpensedefaultJSONResponse{
				Body:       newErrorResponse(ctx, getCausesFromFieldErrors(*errList)),
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		return nil, err
	}
}

func convertUserBalances(userBalances []handlers.ExpenseUserBalance) []v1.ExpenseUserBalance {
	balances := make([]v1.ExpenseUserBalance, 0, len(userBalances))
	for _, ub := range userBalances {
		balances = append(balances, v1.ExpenseUserBalance{
			UserID:  ub.UserId,
			Balance: ub.Balance,
		})
	}
	return balances
}
