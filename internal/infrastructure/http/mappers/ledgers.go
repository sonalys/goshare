package mappers

import (
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func FromLedgersToLedgers(ledgers []domain.Ledger) []server.Ledger {
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
