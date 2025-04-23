package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) AddParticipants(ctx context.Context, ledgerID, userID v1.ID, ids ...v1.ID) error {
	ledgerUUID := convertUUID(ledgerID)

	return mapLedgerError(r.client.transaction(ctx, func(tx pgx.Tx) error {
		query := r.client.queries().WithTx(tx)

		if err := query.LockLedgerForUpdate(ctx, ledgerUUID); err != nil {
			return fmt.Errorf("locking ledger for update: %w", err)
		}

		usersCount, err := query.CountLedgerUsers(ctx, ledgerUUID)
		if err != nil {
			return fmt.Errorf("counting ledger participants: %w", err)
		}

		if usersCount+int64(len(ids)) > v1.LedgerMaxUsers {
			return v1.ErrLedgerMaxUsers
		}

		for _, invitedUserID := range ids {
			addReq := sqlc.AddUserToLedgerParams{
				ID:        convertUUID(v1.NewID()),
				LedgerID:  ledgerUUID,
				UserID:    convertUUID(invitedUserID),
				CreatedAt: convertTime(time.Now()),
				CreatedBy: convertUUID(userID),
			}

			switch err := query.AddUserToLedger(ctx, addReq); {
			case err == nil:
				continue
			case isViolatingConstraint(err, constraintLedgerUniqueParticipant):
				return v1.FieldError{
					Field: "user_id",
					Cause: fmt.Errorf("user '%s' is already a participant of the ledger '%s'", invitedUserID, ledgerID),
				}
			default:
				return fmt.Errorf("adding participant: %w", err)
			}
		}

		return nil
	}))
}

func (r *LedgerRepository) GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error) {
	participants, err := r.client.queries().GetLedgerParticipants(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.LedgerParticipant, 0, len(participants))
	for _, participant := range participants {
		result = append(result, *mappers.NewLedgerParticipant(&participant))
	}
	return result, nil
}
