package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) LedgerParticipantList(ctx context.Context, params handlers.LedgerParticipantListParams) (r *handlers.LedgerParticipantListOK, _ error) {
	balances, err := a.dependencies.LedgerController.GetParticipants(ctx, v1.ConvertID(params.LedgerID))
	if err != nil {
		return nil, err
	}

	return &handlers.LedgerParticipantListOK{
		Participants: mapLedgerParticipantBalanceToResponseObject(balances),
	}, nil
}

func mapLedgerParticipantBalanceToResponseObject(balance []v1.LedgerParticipant) []handlers.LedgerParticipant {
	var balances []handlers.LedgerParticipant
	for _, b := range balance {
		balances = append(balances, handlers.LedgerParticipant{
			ID:        b.ID.UUID(),
			UserID:    b.UserID.UUID(),
			CreatedAt: b.CreatedAt,
			CreatedBy: b.CreatedBy.UUID(),
			Balance:   b.Balance,
		})
	}

	return balances
}
