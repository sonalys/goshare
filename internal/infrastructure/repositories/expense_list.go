package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
	v1 "github.com/sonalys/goshare/pkg/v1"
)

func (r *ExpenseRepository) ListByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error) {
	expenses, err := r.conn.Queries().GetLedgerExpenses(ctx, sqlcgen.GetLedgerExpensesParams{
		LedgerID:  ledgerID,
		Limit:     limit,
		CreatedAt: postgres.ConvertTime(cursor),
	})
	if err != nil {
		return nil, fmt.Errorf("getting ledger expenses: %w", err)
	}

	result := make([]v1.LedgerExpenseSummary, 0, len(expenses))
	for _, expense := range expenses {
		result = append(result, *mappers.NewLedgerExpenseSummary(&expense))
	}

	return result, nil
}
