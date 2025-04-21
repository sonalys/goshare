package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) ListLedgers(ctx context.Context) (r *handlers.ListLedgersOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ledgers, err := a.dependencies.LedgerController.GetByUser(ctx, identity.UserID)
	if err != nil {
		return nil, err
	}

	return &handlers.ListLedgersOK{
		Ledgers: convertLedgers(ledgers),
	}, nil
}

func convertLedgers(ledgers []v1.Ledger) []handlers.Ledger {
	result := make([]handlers.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		result = append(result, handlers.Ledger{
			ID:        ledger.ID.UUID(),
			Name:      ledger.Name,
			CreatedAt: ledger.CreatedAt,
			CreatedBy: ledger.CreatedBy.UUID(),
		})
	}
	return result
}
