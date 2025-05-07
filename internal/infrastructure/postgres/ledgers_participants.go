package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	"github.com/sonalys/kset"
)

func (r *LedgerRepository) Update(ctx context.Context, ledgerID domain.ID, updateFn func(*domain.Ledger) error) error {
	ledgerUUID := convertID(ledgerID)

	return r.transaction(ctx, func(query *sqlc.Queries) error {
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

		updateLedgerParams := sqlc.UpdateLedgerParams{
			ID:   ledgerUUID,
			Name: ledger.Name,
		}
		if err = query.UpdateLedger(ctx, updateLedgerParams); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		currentParticipants := kset.HashMapKey(
			kset.Select(
				func(p sqlc.LedgerParticipant) uuid.UUID { return p.ID.Bytes },
				participantsModel...,
			)...,
		)

		newParticipants := kset.HashMapKeyValue(
			func(p domain.LedgerParticipant) uuid.UUID { return p.ID.UUID() },
			ledger.Participants...,
		)

		for participant := range newParticipants.Difference(currentParticipants).Iter() {
			addReq := sqlc.AddUserToLedgerParams{
				ID:        convertID(domain.NewID()),
				LedgerID:  ledgerUUID,
				UserID:    convertID(participant.Identity),
				Balance:   participant.Balance,
				CreatedAt: convertTime(participant.CreatedAt),
				CreatedBy: convertID(participant.CreatedBy),
			}

			switch err := query.AddUserToLedger(ctx, addReq); {
			case err == nil:
				continue
			case isViolatingConstraint(err, constraintLedgerUniqueParticipant):
				return domain.FieldError{
					Field: "user_id",
					Cause: fmt.Errorf("user '%s' is already a participant of the ledger '%s'", participant.Identity, ledgerID),
				}
			default:
				return fmt.Errorf("adding participant: %w", err)
			}
		}

		for id := range currentParticipants.Difference(newParticipants).Iter() {
			if err := query.RemoveUserFromLedger(ctx, convertUUID(id)); err != nil {
				return fmt.Errorf("removing participant: %w", err)
			}
		}

		return nil
	})
}

func (r *LedgerRepository) GetParticipants(ctx context.Context, ledgerID domain.ID) ([]domain.LedgerParticipant, error) {
	participants, err := r.client.queries().GetLedgerParticipants(ctx, convertID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]domain.LedgerParticipant, 0, len(participants))
	for _, participant := range participants {
		result = append(result, *mappers.NewLedgerParticipant(&participant))
	}
	return result, nil
}
