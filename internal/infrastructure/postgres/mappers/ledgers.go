package mappers

import (
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func NewLedger(ledger *sqlc.Ledger, participants []sqlc.LedgerParticipant) *v1.Ledger {
	ledgerParticipants := make([]v1.LedgerParticipant, 0, len(participants))

	for _, participant := range participants {
		ledgerParticipants = append(ledgerParticipants, v1.LedgerParticipant{
			ID:        newUUID(participant.ID),
			UserID:    newUUID(participant.UserID),
			Balance:   participant.Balance,
			CreatedAt: participant.CreatedAt.Time,
			CreatedBy: newUUID(participant.CreatedBy),
		})
	}

	return &v1.Ledger{
		ID:           newUUID(ledger.ID),
		Name:         ledger.Name,
		Participants: ledgerParticipants,
		CreatedAt:    ledger.CreatedAt.Time,
		CreatedBy:    newUUID(ledger.CreatedBy),
	}
}
