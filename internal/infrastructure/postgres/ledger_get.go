package postgres

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
)

func (r *LedgerRepository) Get(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
	ledger, err := r.client.queries().GetLedgerById(ctx, id)
	if err != nil {
		return nil, ledgerError(err)
	}

	members, err := r.client.queries().GetLedgerMembers(ctx, id)
	if err != nil {
		return nil, ledgerError(err)
	}

	return mappers.NewLedger(&ledger, members), nil
}
