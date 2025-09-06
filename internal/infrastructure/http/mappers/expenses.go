package mappers

import (
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func LedgerExpenseSummaryToExpenseSummary(expenses []application.LedgerExpenseSummary) []server.ExpenseSummary {
	expensesResponse := make([]server.ExpenseSummary, 0, len(expenses))

	for _, e := range expenses {
		expensesResponse = append(expensesResponse, server.ExpenseSummary{
			ID:          e.ID.UUID(),
			Amount:      e.Amount,
			Name:        e.Name,
			ExpenseDate: e.ExpenseDate,
			CreatedAt:   e.CreatedAt,
			CreatedBy:   e.CreatedBy.UUID(),
			UpdatedAt:   e.UpdatedAt,
			UpdatedBy:   e.UpdatedBy.UUID(),
		})
	}

	return expensesResponse
}
