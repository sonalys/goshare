package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) ListLedgers(ctx context.Context, request handlers.ListLedgersRequestObject) (handlers.ListLedgersResponseObject, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ledgers, err := a.dependencies.UserLedgerLister.GetByUser(ctx, identity.UserID)
	if err != nil {
		return nil, err
	}

	return handlers.ListLedgers200JSONResponse{
		Ledgers: convertLedgers(ledgers),
	}, nil
}

func convertLedgers(ledgers []v1.Ledger) []handlers.Ledger {
	result := make([]handlers.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		result = append(result, handlers.Ledger{
			Id:        ledger.ID,
			Name:      ledger.Name,
			CreatedAt: ledger.CreatedAt,
			CreatedBy: ledger.CreatedBy,
		})
	}
	return result
}
