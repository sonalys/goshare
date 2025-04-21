package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	"github.com/sonalys/goshare/internal/pkg/pointers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) CreateExpense(ctx context.Context, req *handlers.CreateExpenseReq, params handlers.CreateExpenseParams) (r *handlers.CreateExpenseOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var categoryID *v1.ID
	if id, ok := req.CategoryID.Get(); ok {
		categoryID = pointers.New(v1.ConvertID(id))
	}

	apiReq := ledgers.CreateExpenseRequest{
		UserID:       identity.UserID,
		LedgerID:     v1.ConvertID(params.LedgerID),
		CategoryID:   categoryID,
		Amount:       req.Amount,
		Name:         req.Name,
		ExpenseDate:  req.ExpenseDate,
		UserBalances: convertUserBalances(req.UserBalances),
	}

	switch resp, err := a.dependencies.ExpenseCreater.CreateExpense(ctx, apiReq); {
	case err == nil:
		return &handlers.CreateExpenseOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}

func convertUserBalances(userBalances []handlers.ExpenseUserBalance) []v1.ExpenseUserBalance {
	balances := make([]v1.ExpenseUserBalance, 0, len(userBalances))
	for _, ub := range userBalances {
		balances = append(balances, v1.ExpenseUserBalance{
			UserID:  v1.ConvertID(ub.UserID),
			Balance: ub.Balance,
		})
	}
	return balances
}
