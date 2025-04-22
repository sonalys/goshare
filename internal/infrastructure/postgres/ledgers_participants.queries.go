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

	return nil
}
