package mappers

import (
	domain "github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func NewLedger(ledger *sqlc.Ledger, participants []sqlc.LedgerParticipant) *domain.Ledger {
	ledgerParticipants := make([]domain.LedgerMember, 0, len(participants))

	for _, participant := range participants {
		ledgerParticipants = append(ledgerParticipants, domain.LedgerMember{
			ID:        newUUID(participant.ID),
			Identity:  newUUID(participant.UserID),
			Balance:   participant.Balance,
			CreatedAt: participant.CreatedAt.Time,
			CreatedBy: newUUID(participant.CreatedBy),
		})
	}

	return &domain.Ledger{
		ID:        newUUID(ledger.ID),
		Name:      ledger.Name,
		Members:   ledgerParticipants,
		CreatedAt: ledger.CreatedAt.Time,
		CreatedBy: newUUID(ledger.CreatedBy),
	}
}
