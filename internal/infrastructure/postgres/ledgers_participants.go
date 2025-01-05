package postgres

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) AddParticipant(ctx context.Context, ledgerID, userID, invitedUserID v1.ID) error {
	return mapLedgerError(r.client.transaction(ctx, func(tx *sqlc.Queries) error {
		return addLedgerParticipant(ctx, tx, ledgerID, userID, invitedUserID)
	}))
}

func (r *LedgerRepository) GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error) {
	participants, err := r.client.queries().GetLedgerParticipants(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.LedgerParticipant, 0, len(participants))
	for _, user := range participants {
		result = append(result, *newLedgerParticipant(&user))
	}
	return result, nil
}

func newLedgerParticipant(user *sqlc.LedgerParticipant) *v1.LedgerParticipant {
	return &v1.LedgerParticipant{
		ID:        newUUID(user.ID),
		LedgerID:  newUUID(user.LedgerID),
		UserID:    newUUID(user.UserID),
		CreatedAt: user.CreatedAt.Time,
		CreatedBy: newUUID(user.CreatedBy),
	}
}
