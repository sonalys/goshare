package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) AddParticipants(ctx context.Context, ledgerID v1.ID, updateFn func(*v1.Ledger) error) error {
	ledgerUUID := convertUUID(ledgerID)

	return mapLedgerError(r.client.transaction(ctx, func(tx pgx.Tx) error {
		query := r.client.queries().WithTx(tx)

		if err := query.LockLedgerForUpdate(ctx, ledgerUUID); err != nil {
			return fmt.Errorf("locking ledger for update: %w", err)
		}

		ledgerModel, err := query.FindLedgerById(ctx, ledgerUUID)
		if err != nil {
			return fmt.Errorf("getting ledger: %w", err)
		}

		participantsModel, err := query.GetLedgerParticipants(ctx, ledgerUUID)
		if err != nil {
			return fmt.Errorf("getting ledger participants: %w", err)
		}

		ledger := mappers.NewLedger(&ledgerModel, participantsModel)

		if err := updateFn(ledger); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		participantsToRemove := make([]sqlc.LedgerParticipant, 0, len(ledger.Participants))
		participantsToAdd := make([]v1.LedgerParticipant, 0, len(ledger.Participants))

	outer:
		for _, participant := range ledger.Participants {
			for _, existingParticipant := range participantsModel {
				if participant.UserID.UUID() == existingParticipant.UserID.Bytes {
					continue outer
				}

				participantsToAdd = append(participantsToAdd, participant)
			}
		}

		for _, existingParticipant := range participantsModel {
			for _, participant := range ledger.Participants {
				if participant.UserID.UUID() == existingParticipant.UserID.Bytes {
					continue outer
				}
			}
			participantsToRemove = append(participantsToRemove, existingParticipant)
		}

		for _, participant := range participantsToAdd {
			addReq := sqlc.AddUserToLedgerParams{
				ID:        convertUUID(v1.NewID()),
				LedgerID:  ledgerUUID,
				UserID:    convertUUID(participant.UserID),
				CreatedAt: convertTime(participant.CreatedAt),
				CreatedBy: convertUUID(participant.CreatedBy),
			}

			switch err := query.AddUserToLedger(ctx, addReq); {
			case err == nil:
				continue
			case isViolatingConstraint(err, constraintLedgerUniqueParticipant):
				return v1.FieldError{
					Field: "user_id",
					Cause: fmt.Errorf("user '%s' is already a participant of the ledger '%s'", participant.UserID, ledgerID),
				}
			default:
				return fmt.Errorf("adding participant: %w", err)
			}
		}

		for _, participant := range participantsToRemove {
			if err := query.RemoveUserFromLedger(ctx, participant.ID); err != nil {
				return fmt.Errorf("removing participant: %w", err)
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
