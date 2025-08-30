package postgres

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *LedgerRepository) Create(ctx context.Context, ledger *domain.Ledger) error {
	return r.transaction(ctx, func(query *sqlcgen.Queries) error {
		createLedgerReq := sqlcgen.CreateLedgerParams{
			ID:        ledger.ID,
			Name:      ledger.Name,
			CreatedAt: convertTime(ledger.CreatedAt),
			CreatedBy: ledger.CreatedBy,
		}

		if err := query.CreateLedger(ctx, createLedgerReq); err != nil {
			return fmt.Errorf("failed to create ledger: %w", err)
		}

		for id, member := range ledger.Members {
			addReq := sqlcgen.CreateLedgerMemberParams{
				UserID:    id,
				LedgerID:  createLedgerReq.ID,
				CreatedAt: convertTime(member.CreatedAt),
				CreatedBy: member.CreatedBy,
				Balance:   member.Balance,
			}

			if err := query.CreateLedgerMember(ctx, addReq); err != nil {
				return fmt.Errorf("failed to add user to ledger: %w", err)
			}
		}

		return nil
	})
}
