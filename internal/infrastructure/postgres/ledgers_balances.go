package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) GetLedgerBalance(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerParticipantBalance, error) {
	balances, err := r.client.queries().GetLedgerBalances(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.LedgerParticipantBalance, 0, len(balances))
	for _, balance := range balances {
		result = append(result, *newLedgerParticipantBalance(&balance))
	}
	return result, nil
}

func newLedgerParticipantBalance(balance *sqlc.LedgerParticipantBalance) *v1.LedgerParticipantBalance {
	return &v1.LedgerParticipantBalance{
		ID:            newUUID(balance.ID),
		LedgerID:      newUUID(balance.LedgerID),
		UserID:        newUUID(balance.UserID),
		LastTimestamp: balance.LastTimestamp.Time,
		Balance:       balance.Balance,
	}
}
