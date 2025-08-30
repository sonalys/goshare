package postgres

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
)

func (r *LedgerRepository) ListByUser(ctx context.Context, userID domain.ID) ([]domain.Ledger, error) {
	ledgers, err := r.client.queries().GetUserLedgers(ctx, userID)
	if err != nil {
		return nil, ledgerError(err)
	}

	result := make([]domain.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		members, err := r.client.queries().GetLedgerMembers(ctx, ledger.ID)
		if err != nil {
			return nil, ledgerError(err)
		}
		result = append(result, *mappers.NewLedger(&ledger, members))
	}
	return result, nil
}
