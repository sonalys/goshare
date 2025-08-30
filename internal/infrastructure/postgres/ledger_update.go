package postgres

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *LedgerRepository) Update(ctx context.Context, ledger *domain.Ledger) error {
	return r.transaction(ctx, func(conn connection) error {
		query := conn.queries()

		updateLedgerParams := sqlcgen.UpdateLedgerParams{
			ID:   ledger.ID,
			Name: ledger.Name,
		}
		if err := query.UpdateLedger(ctx, updateLedgerParams); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		memberIDs := slices.Collect(maps.Keys(ledger.Members))

		if err := query.DeleteMembersNotIn(ctx, memberIDs); err != nil {
			return fmt.Errorf("deleting old members: %w", err)
		}

		for id, member := range ledger.Members {
			err := query.CreateLedgerMember(ctx, sqlcgen.CreateLedgerMemberParams{
				LedgerID:  ledger.ID,
				UserID:    id,
				CreatedAt: convertTime(member.CreatedAt),
				CreatedBy: member.CreatedBy,
				Balance:   member.Balance,
			})
			if err != nil {
				return fmt.Errorf("saving ledger member: %w", err)
			}
		}

		return nil
	})
}
