package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *Repository) ListByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]application.LedgerExpenseSummary, error) {
	expenses, err := r.conn.Queries().GetLedgerExpenses(ctx, sqlcgen.GetLedgerExpensesParams{
		LedgerID:  ledgerID,
		Limit:     limit,
		CreatedAt: postgres.ConvertTime(cursor),
	})
	if err != nil {
		return nil, fmt.Errorf("getting ledger expenses: %w", expenseError(err))
	}

	result := make([]application.LedgerExpenseSummary, 0, len(expenses))
	for _, expense := range expenses {
		result = append(result, *toLedgerExpenseSummary(&expense))
	}

	return result, nil
}
