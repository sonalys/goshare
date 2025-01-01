package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/queries"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) AddParticipant(ctx context.Context, ledgerID, userID, invitedUserID uuid.UUID) error {
	return mapLedgerError(r.client.queries().AddUserToLedger(ctx, queries.AddUserToLedgerParams{
		LedgerID:  convertUUID(ledgerID),
		UserID:    convertUUID(invitedUserID),
		ID:        convertUUID(uuid.New()),
		CreatedAt: convertTime(time.Now()),
		CreatedBy: convertUUID(userID),
	}))
}

func (r *LedgerRepository) GetParticipants(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerParticipant, error) {
	users, err := r.client.queries().GetLedgerParticipants(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.LedgerParticipant, 0, len(users))
	for _, user := range users {
		result = append(result, *newLedgerParticipant(&user))
	}
	return result, nil
}

func newLedgerParticipant(user *queries.LedgerParticipant) *v1.LedgerParticipant {
	return &v1.LedgerParticipant{
		ID:        newUUID(user.ID),
		LedgerID:  newUUID(user.LedgerID),
		UserID:    newUUID(user.UserID),
		CreatedAt: user.CreatedAt.Time,
		CreatedBy: newUUID(user.CreatedBy),
	}
}
