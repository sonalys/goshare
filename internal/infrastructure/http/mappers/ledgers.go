package mappers

import (
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func LedgersToLedgers(ledgers []domain.Ledger) []server.Ledger {
	result := make([]server.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		result = append(result, server.Ledger{
			ID:        ledger.ID.UUID(),
			Name:      ledger.Name,
			CreatedAt: ledger.CreatedAt,
			CreatedBy: ledger.CreatedBy.UUID(),
		})
	}

	return result
}

func LedgerMemberToLedgerMember(members map[domain.ID]*domain.LedgerMember) []server.LedgerMember {
	balances := make([]server.LedgerMember, 0, len(members))
	for id, member := range members {
		balances = append(balances, server.LedgerMember{
			UserID:    id.UUID(),
			CreatedAt: member.CreatedAt,
			CreatedBy: member.CreatedBy.UUID(),
			Balance:   member.Balance,
		})
	}

	return balances
}
