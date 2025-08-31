package ledger

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
)

func (r *Repository) ListByUser(ctx context.Context, userID domain.ID) ([]domain.Ledger, error) {
	ledgers, err := r.client.Queries().GetUserLedgers(ctx, userID)
	if err != nil {
		return nil, ledgerError(err)
	}

	result := make([]domain.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		members, err := r.client.Queries().GetLedgerMembers(ctx, ledger.ID)
		if err != nil {
			return nil, ledgerError(err)
		}
		result = append(result, *toLedger(&ledger, members))
	}

	return result, nil
}
