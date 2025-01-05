package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func addLedgerParticipant(ctx context.Context, tx *sqlc.Queries, ledgerID, userID, invitedUserID uuid.UUID) error {
	addReq := sqlc.AddUserToLedgerParams{
		LedgerID:  convertUUID(ledgerID),
		UserID:    convertUUID(invitedUserID),
		ID:        convertUUID(uuid.New()),
		CreatedAt: convertTime(time.Now()),
		CreatedBy: convertUUID(userID),
	}

	if err := tx.AddUserToLedger(ctx, addReq); err != nil {
		return fmt.Errorf("failed to add user to ledger: %w", err)
	}

	if err := createLedgerParticipantBalance(ctx, tx, ledgerID, invitedUserID); err != nil {
		return fmt.Errorf("failed to create ledger participant balance: %w", err)
	}

	return nil
}

func createLedgerParticipantBalance(ctx context.Context, tx *sqlc.Queries, ledgerID, userID uuid.UUID) error {
	return tx.CreateLedgerParticipantBalance(ctx, sqlc.CreateLedgerParticipantBalanceParams{
		ID:            convertUUID(uuid.New()),
		LedgerID:      convertUUID(ledgerID),
		UserID:        convertUUID(userID),
		LastTimestamp: convertTime(time.Now()),
		Balance:       0,
	})
}

func getOldestTimestamp(balances []sqlc.LedgerParticipantBalance) time.Time {
	var oldestTimestamp time.Time

	for _, balance := range balances {
		if balance.LastTimestamp.Time.Before(oldestTimestamp) {
			oldestTimestamp = balance.LastTimestamp.Time
		}
	}

	return oldestTimestamp
}

func updateLedgerParticipantsBalance(ctx context.Context, tx *sqlc.Queries, ledgerID uuid.UUID) error {
	balances, err := tx.GetLedgerBalances(ctx, convertUUID(ledgerID))
	if err != nil {
		return fmt.Errorf("failed to get ledger participants balances: %w", err)
	}

	oldestTimestamp := getOldestTimestamp(balances)

	// records are sorted by created_at in ascending order.
	records, err := tx.GetLedgerRecordsFromTimestamp(ctx, sqlc.GetLedgerRecordsFromTimestampParams{
		LedgerID:  convertUUID(ledgerID),
		CreatedAt: convertTime(oldestTimestamp),
	})
	if err != nil {
		return fmt.Errorf("failed to get ledger records from timestamp: %w", err)
	}

	for _, record := range records {
		for i := range balances {
			balance := &balances[i]
			if record.UserID == balance.UserID && record.CreatedAt.Time.After(balance.LastTimestamp.Time) {
				balance.Balance += record.Amount
				balance.LastTimestamp = record.CreatedAt
			}
		}
	}

	for i := range balances {
		balance := &balances[i]
		updateReq := sqlc.UpdateLedgerParticipantBalanceParams{
			UserID:        balance.UserID,
			LedgerID:      balance.LedgerID,
			LastTimestamp: balance.LastTimestamp,
			Balance:       balance.Balance,
		}
		if err := tx.UpdateLedgerParticipantBalance(ctx, updateReq); err != nil {
			return fmt.Errorf("failed to update ledger participant balance: %w", err)
		}
	}

	return nil
}
