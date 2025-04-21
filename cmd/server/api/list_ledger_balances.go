package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) ListLedgerBalances(ctx context.Context, params handlers.ListLedgerBalancesParams) (r *handlers.ListLedgerBalancesOK, _ error) {
	balances, err := a.dependencies.LedgerController.GetBalances(ctx, v1.ConvertID(params.LedgerID))
	if err != nil {
		return nil, err
	}

	return &handlers.ListLedgerBalancesOK{
		Balances: mapLedgerParticipantBalanceToResponseObject(balances),
	}, nil
}

func mapLedgerParticipantBalanceToResponseObject(balance []v1.LedgerParticipantBalance) []handlers.LedgerParticipantBalance {
	var balances []handlers.LedgerParticipantBalance
	for _, b := range balance {
		balances = append(balances, handlers.LedgerParticipantBalance{
			UserID:  b.UserID.UUID(),
			Balance: b.Balance,
		})
	}

	return balances
}
