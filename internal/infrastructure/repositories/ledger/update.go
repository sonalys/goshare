package ledger

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func updateLedgerParams(ledger *domain.Ledger) sqlcgen.UpdateLedgerParams {
	return sqlcgen.UpdateLedgerParams{
		ID:   ledger.ID,
		Name: ledger.Name,
	}
}

func (r *Repository) Update(ctx context.Context, ledger *domain.Ledger) error {
	return r.transaction(ctx, func(conn postgres.Connection) error {
		query := conn.Queries()

		if err := query.UpdateLedger(ctx, updateLedgerParams(ledger)); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		memberIDs := slices.Collect(maps.Keys(ledger.Members))

		if err := query.DeleteMembersNotIn(ctx, memberIDs); err != nil {
			return fmt.Errorf("deleting old members: %w", err)
		}

		for id, member := range ledger.Members {
			err := query.CreateLedgerMember(ctx, createLedgerMemberParams(ledger.ID, id, member))
			if err != nil {
				return fmt.Errorf("saving ledger member: %w", err)
			}
		}

		return nil
	})
}
