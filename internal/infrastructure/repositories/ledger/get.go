package ledger

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
)

func (r *Repository) Get(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
	ledger, err := r.client.Queries().GetLedgerById(ctx, id)
	if err != nil {
		return nil, ledgerError(err)
	}

	members, err := r.client.Queries().GetLedgerMembers(ctx, id)
	if err != nil {
		return nil, ledgerError(err)
	}

	return toLedger(&ledger, members), nil
}
