package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) ListLedgerBalances(ctx context.Context, request handlers.ListLedgerBalancesRequestObject) (handlers.ListLedgerBalancesResponseObject, error) {
	balances, err := a.dependencies.LedgerBalancesLister.GetBalances(ctx, request.LedgerID)
	if err != nil {
		return nil, err
	}

	return handlers.ListLedgerBalances200JSONResponse{
		Balances: mapLedgerParticipantBalanceToResponseObject(balances),
	}, nil
}

func mapLedgerParticipantBalanceToResponseObject(balance []v1.LedgerParticipantBalance) []handlers.LedgerParticipantBalance {
	var balances []handlers.LedgerParticipantBalance
	for _, b := range balance {
		balances = append(balances, handlers.LedgerParticipantBalance{
			UserId:  b.UserID,
			Balance: float32(b.Balance),
		})
	}

	return balances
}
