package postgres

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) GetLedgerBalance(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerParticipantBalance, error) {
	balances, err := r.client.queries().GetLedgerParticipantsWithBalance(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.LedgerParticipantBalance, 0, len(balances))
	for _, balance := range balances {
		result = append(result, *newLedgerParticipantBalance(&balance))
	}

	go func() {
		ctx := context.Background()
		for i := range balances {
			balance := &balances[i]
			err := r.client.queries().UpsertLedgerParticipantBalance(ctx, sqlc.UpsertLedgerParticipantBalanceParams{
				ID:            convertUUID(uuid.New()),
				LedgerID:      balance.LedgerID,
				UserID:        balance.UserID,
				LastTimestamp: balance.LastTimestamp,
				Balance:       balance.Balance,
			})
			if err != nil {
				slog.Warn("failed to update ledger participant balance",
					slog.String("user_id", newUUID(balance.UserID).String()),
					slog.String("ledger_id", newUUID(balance.LedgerID).String()),
					slog.Any("error", err),
				)
			}
		}
	}()
	return result, nil
}

func newLedgerParticipantBalance(balance *sqlc.GetLedgerParticipantsWithBalanceRow) *v1.LedgerParticipantBalance {
	return &v1.LedgerParticipantBalance{
		LedgerID: newUUID(balance.LedgerID),
		UserID:   newUUID(balance.UserID),
		Balance:  balance.Balance,
	}
}
