package api

import (
	"context"
	"time"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) LedgerExpenseCreate(ctx context.Context, req *handlers.Expense, params handlers.LedgerExpenseCreateParams) (r *handlers.LedgerExpenseCreateOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	apiReq := ledgers.CreateExpenseRequest{
		UserID:      identity.UserID,
		LedgerID:    v1.ConvertID(params.LedgerID),
		Name:        req.Name,
		ExpenseDate: req.ExpenseDate,
		Records:     convertUserBalances(identity.UserID, req.Records),
	}

	switch resp, err := a.dependencies.LedgerController.CreateExpense(ctx, apiReq); err {
	case nil:
		return &handlers.LedgerExpenseCreateOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}

func convertUserBalances(identity v1.ID, userBalances []handlers.ExpenseRecord) []v1.Record {
	balances := make([]v1.Record, 0, len(userBalances))
	for _, ub := range userBalances {
		balances = append(balances, v1.Record{
			ID:        v1.NewID(),
			Type:      v1.NewRecordType(string(ub.Type)),
			Amount:    ub.Amount,
			From:      v1.ConvertID(ub.FromUserID),
			To:        v1.ConvertID(ub.ToUserID),
			CreatedAt: time.Now(),
			CreatedBy: identity,
			UpdatedAt: time.Now(),
			UpdatedBy: identity,
		})
	}
	return balances
}
