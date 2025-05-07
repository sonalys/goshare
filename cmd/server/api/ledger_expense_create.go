package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerExpenseCreate(ctx context.Context, req *handlers.Expense, params handlers.LedgerExpenseCreateParams) (r *handlers.LedgerExpenseCreateOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	apiReq := controllers.CreateExpenseRequest{
		Identity:    identity.UserID,
		LedgerID:    domain.ConvertID(params.LedgerID),
		Name:        req.Name,
		ExpenseDate: req.ExpenseDate,
		Records:     convertUserBalances(req.Records),
	}

	switch resp, err := a.Ledgers.CreateExpense(ctx, apiReq); err {
	case nil:
		return &handlers.LedgerExpenseCreateOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}

func convertUserBalances(userBalances []handlers.ExpenseRecord) []domain.NewRecord {
	balances := make([]domain.NewRecord, 0, len(userBalances))
	for _, ub := range userBalances {
		balances = append(balances, domain.NewRecord{
			Type:   domain.NewRecordType(string(ub.Type)),
			Amount: ub.Amount,
			From:   domain.ConvertID(ub.FromUserID),
			To:     domain.ConvertID(ub.ToUserID),
		})
	}
	return balances
}
