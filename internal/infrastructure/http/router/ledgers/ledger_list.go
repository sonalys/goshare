package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerList(ctx context.Context) (*server.LedgerListOK, error) {
	identity, err := middlewares.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ledgers, err := a.UserController.ListByUser(ctx, identity.UserID)
	if err != nil {
		return nil, err
	}

	return &server.LedgerListOK{
		Ledgers: convertLedgers(ledgers),
	}, nil
}

func convertLedgers(ledgers []domain.Ledger) []server.Ledger {
	result := make([]server.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		result = append(result, server.Ledger{
			ID:        ledger.ID.UUID(),
			Name:      ledger.Name,
			CreatedAt: ledger.CreatedAt,
			CreatedBy: ledger.CreatedBy.UUID(),
		})
	}

	return result
}
