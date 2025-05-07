package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerList(ctx context.Context) (*handlers.LedgerListOK, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ledgers, err := a.Ledgers.GetByUser(ctx, identity.UserID)
	if err != nil {
		return nil, err
	}

	return &handlers.LedgerListOK{
		Ledgers: convertLedgers(ledgers),
	}, nil
}

func convertLedgers(ledgers []domain.Ledger) []handlers.Ledger {
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
