package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerMemberList(ctx context.Context, params handlers.LedgerMemberListParams) (r *handlers.LedgerMemberListOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	members, err := a.UserController.ListMembers(ctx, usercontroller.ListMembersRequest{
		Actor:    identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
	})
	if err != nil {
		return nil, err
	}

	return &handlers.LedgerMemberListOK{
		Members: mapLedgerMemberBalanceToResponseObject(members),
	}, nil
}

func mapLedgerMemberBalanceToResponseObject(members map[domain.ID]*domain.LedgerMember) []handlers.LedgerMember {
	var balances []handlers.LedgerMember
	for id, member := range members {
		balances = append(balances, handlers.LedgerMember{
			UserID:    id.UUID(),
			CreatedAt: member.CreatedAt,
			CreatedBy: member.CreatedBy.UUID(),
			Balance:   member.Balance,
		})
	}

	return balances
}
