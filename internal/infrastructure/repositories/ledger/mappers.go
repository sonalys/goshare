package ledger

import (
	domain "github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func toLedger(ledger *sqlcgen.Ledger, members []sqlcgen.LedgerMember) *domain.Ledger {
	ledgerMembers := make(map[domain.ID]*domain.LedgerMember, len(members))

	for _, member := range members {
		ledgerMembers[member.UserID] = &domain.LedgerMember{
			Balance:   member.Balance,
			CreatedAt: member.CreatedAt.Time,
			CreatedBy: member.CreatedBy,
		}
	}

	return &domain.Ledger{
		ID:        ledger.ID,
		Name:      ledger.Name,
		Members:   ledgerMembers,
		CreatedAt: ledger.CreatedAt.Time,
		CreatedBy: ledger.CreatedBy,
	}
}
