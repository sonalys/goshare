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
	go r.updateBalanceSnapshot(ledgerID, balances)

	result := make([]v1.LedgerParticipantBalance, 0, len(balances))
	for _, balance := range balances {
		result = append(result, *newLedgerParticipantBalance(&balance))
	}
	return result, nil
}

func (r *LedgerRepository) updateBalanceSnapshot(ledgerID uuid.UUID, balances []sqlc.GetLedgerParticipantsWithBalanceRow) {
	ctx := context.Background()

	fields := []any{
		slog.String("ledger_id", ledgerID.String()),
	}

	slog.InfoContext(ctx, "starting balance snapshot update", fields...)

	var updateCount int
	for i := range balances {
		balance := &balances[i]
		if balance.LastTimestamp.Time.IsZero() {
			continue
		}
		err := r.client.queries().UpsertLedgerParticipantBalance(ctx, sqlc.UpsertLedgerParticipantBalanceParams{
			ID:            convertUUID(uuid.New()),
			LedgerID:      balance.LedgerID,
			UserID:        balance.UserID,
			LastTimestamp: balance.LastTimestamp,
			Balance:       balance.Balance,
		})
		if err != nil {
			slog.Warn("failed to update ledger participant balance",
				append(fields,
					slog.String("user_id", newUUID(balance.UserID).String()),
					slog.Any("error", err),
				)...,
			)
			continue
		}
		updateCount++
		slog.InfoContext(ctx, "user balance snapshot updated",
			append(fields,
				slog.String("user_id", newUUID(balance.UserID).String()),
			)...,
		)
	}

	slog.InfoContext(ctx, "balance snapshot update finished",
		append(fields,
			slog.Int("updated_entries", updateCount),
		)...,
	)
}

func newLedgerParticipantBalance(balance *sqlc.GetLedgerParticipantsWithBalanceRow) *v1.LedgerParticipantBalance {
	return &v1.LedgerParticipantBalance{
		LedgerID: newUUID(balance.LedgerID),
		UserID:   newUUID(balance.UserID),
		Balance:  balance.Balance,
	}
}
