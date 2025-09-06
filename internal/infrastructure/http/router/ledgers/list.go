package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/http/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerList(ctx context.Context) (*server.LedgerListOK, error) {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ledgers, err := a.Ledgers().ListByUser(ctx, identity.UserID)
	if err != nil {
		return nil, err
	}

	return &server.LedgerListOK{
		Ledgers: mappers.LedgersToLedgers(ledgers),
	}, nil
}
