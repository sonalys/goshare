package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerMemberList(ctx context.Context, params server.LedgerMemberListParams) (r *server.LedgerMemberListOK, _ error) {
	identity, err := middlewares.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ledger, err := a.UserController.Ledgers().Get(ctx, usercontroller.GetLedgerRequest{
		ActorID:  identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
	})
	if err != nil {
		return nil, err
	}

	return &server.LedgerMemberListOK{
		Members: mapLedgerMemberBalanceToResponseObject(ledger.Members),
	}, nil
}

func mapLedgerMemberBalanceToResponseObject(members map[domain.ID]*domain.LedgerMember) []server.LedgerMember {
	balances := make([]server.LedgerMember, 0, len(members))
	for id, member := range members {
		balances = append(balances, server.LedgerMember{
			UserID:    id.UUID(),
			CreatedAt: member.CreatedAt,
			CreatedBy: member.CreatedBy.UUID(),
			Balance:   member.Balance,
		})
	}

	return balances
}
