package members

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerMemberList(ctx context.Context, params server.LedgerMemberListParams) (r *server.LedgerMemberListOK, _ error) {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ledger, err := a.Ledgers().Get(ctx, usercontroller.GetLedgerRequest{
		ActorID:  identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
	})
	if err != nil {
		return nil, err
	}

	return &server.LedgerMemberListOK{
		Members: mappers.LedgerMemberToLedgerMember(ledger.Members),
	}, nil
}
