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
