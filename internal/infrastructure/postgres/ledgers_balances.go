package postgres

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) GetLedgerBalance(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipantBalance, error) {
	balances, err := r.client.queries().GetLedgerParticipantsWithBalance(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	result := make([]v1.LedgerParticipantBalance, 0, len(balances))
	for _, balance := range balances {
		result = append(result, *newLedgerParticipantBalance(&balance))
	}
	return result, nil
}

func newLedgerParticipantBalance(balance *sqlc.GetLedgerParticipantsWithBalanceRow) *v1.LedgerParticipantBalance {
	return &v1.LedgerParticipantBalance{
		LedgerID: newUUID(balance.LedgerID),
		UserID:   newUUID(balance.UserID),
		Balance:  balance.Balance,
	}
}
