package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerParticipantList(ctx context.Context, params handlers.LedgerParticipantListParams) (r *handlers.LedgerParticipantListOK, _ error) {
	balances, err := a.Ledgers.GetParticipants(ctx, domain.ConvertID(params.LedgerID))
	if err != nil {
		return nil, err
	}

	return &handlers.LedgerParticipantListOK{
		Participants: mapLedgerParticipantBalanceToResponseObject(balances),
	}, nil
}

func mapLedgerParticipantBalanceToResponseObject(balance []domain.LedgerParticipant) []handlers.LedgerParticipant {
	var balances []handlers.LedgerParticipant
	for _, b := range balance {
		balances = append(balances, handlers.LedgerParticipant{
			ID:        b.ID.UUID(),
			UserID:    b.Identity.UUID(),
			CreatedAt: b.CreatedAt,
			CreatedBy: b.CreatedBy.UUID(),
			Balance:   b.Balance,
		})
	}

	return balances
}
