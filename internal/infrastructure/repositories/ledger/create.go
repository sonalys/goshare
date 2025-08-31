package ledger

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func createLedgerParams(ledger *domain.Ledger) sqlcgen.CreateLedgerParams {
	return sqlcgen.CreateLedgerParams{
		ID:        ledger.ID,
		Name:      ledger.Name,
		CreatedAt: postgres.ConvertTime(ledger.CreatedAt),
		CreatedBy: ledger.CreatedBy,
	}
}

func createLedgerMemberParams(ledgerID, memberID domain.ID, member *domain.LedgerMember) sqlcgen.CreateLedgerMemberParams {
	return sqlcgen.CreateLedgerMemberParams{
		UserID:    memberID,
		LedgerID:  ledgerID,
		CreatedAt: postgres.ConvertTime(member.CreatedAt),
		CreatedBy: member.CreatedBy,
		Balance:   member.Balance,
	}
}

func (r *Repository) Create(ctx context.Context, ledger *domain.Ledger) error {
	return r.transaction(ctx, func(conn postgres.Connection) error {
		query := conn.Queries()

		if err := query.CreateLedger(ctx, createLedgerParams(ledger)); err != nil {
			return fmt.Errorf("failed to create ledger: %w", err)
		}

		for id, member := range ledger.Members {
			if err := query.CreateLedgerMember(ctx, createLedgerMemberParams(ledger.ID, id, member)); err != nil {
				return fmt.Errorf("failed to add user to ledger: %w", err)
			}
		}

		return nil
	})
}
