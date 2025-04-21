package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) addLedgerParticipant(ctx context.Context, queries *sqlc.Queries, ledgerID, userID, invitedUserID v1.ID) error {
	addReq := sqlc.AddUserToLedgerParams{
		LedgerID:  convertUUID(ledgerID),
		UserID:    convertUUID(invitedUserID),
		ID:        convertUUID(v1.NewID()),
		CreatedAt: convertTime(time.Now()),
		CreatedBy: convertUUID(userID),
	}

	if err := queries.AddUserToLedger(ctx, addReq); err != nil {
		return fmt.Errorf("failed to add user to ledger: %w", err)
	}

	createLedgerBalanceReq := sqlc.CreateLedgerParticipantBalanceParams{
		ID:            convertUUID(v1.NewID()),
		LedgerID:      convertUUID(ledgerID),
		UserID:        convertUUID(invitedUserID),
		LastTimestamp: convertTime(time.Now()),
		Balance:       0,
	}
	if err := queries.CreateLedgerParticipantBalance(ctx, createLedgerBalanceReq); err != nil {
		return fmt.Errorf("failed to create ledger participant balance: %w", err)
	}

	return nil
}
