package postgres

import (
	"context"
	"fmt"
	"maps"
	"slices"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

type LedgerRepository struct {
	client connection
}

func NewLedgerRepository(client connection) *LedgerRepository {
	return &LedgerRepository{
		client: client,
	}
}

func (r *LedgerRepository) transaction(ctx context.Context, f func(q connection) error) error {
	return mapLedgerError(r.client.transaction(ctx, f))
}

func (r *LedgerRepository) Create(ctx context.Context, ledger *domain.Ledger) error {
	return r.transaction(ctx, func(conn connection) error {
		query := conn.queries()

		createLedgerReq := sqlc.CreateLedgerParams{
			ID:        ledger.ID,
			Name:      ledger.Name,
			CreatedAt: convertTime(ledger.CreatedAt),
			CreatedBy: ledger.CreatedBy,
		}

		if err := query.CreateLedger(ctx, createLedgerReq); err != nil {
			return fmt.Errorf("failed to create ledger: %w", err)
		}

		for id, member := range ledger.Members {
			addReq := sqlc.SaveLedgerMemberParams{
				UserID:    id,
				LedgerID:  createLedgerReq.ID,
				CreatedAt: convertTime(member.CreatedAt),
				CreatedBy: member.CreatedBy,
				Balance:   member.Balance,
			}

			if err := query.SaveLedgerMember(ctx, addReq); err != nil {
				return fmt.Errorf("failed to add user to ledger: %w", err)
			}
		}

		return nil
	})
}

func (r *LedgerRepository) Get(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
	ledger, err := r.client.queries().GetLedgerById(ctx, id)
	if err != nil {
		return nil, mapLedgerError(err)
	}

	members, err := r.client.queries().GetLedgerMembers(ctx, id)
	if err != nil {
		return nil, mapLedgerError(err)
	}

	return mappers.NewLedger(&ledger, members), nil
}

func (r *LedgerRepository) ListByUser(ctx context.Context, userID domain.ID) ([]domain.Ledger, error) {
	ledgers, err := r.client.queries().GetUserLedgers(ctx, userID)
	if err != nil {
		return nil, mapLedgerError(err)
	}

	result := make([]domain.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		members, err := r.client.queries().GetLedgerMembers(ctx, ledger.ID)
		if err != nil {
			return nil, mapLedgerError(err)
		}
		result = append(result, *mappers.NewLedger(&ledger, members))
	}
	return result, nil
}

func (r *LedgerRepository) Update(ctx context.Context, ledger *domain.Ledger) error {
	return r.transaction(ctx, func(conn connection) error {
		query := conn.queries()

		updateLedgerParams := sqlc.UpdateLedgerParams{
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
			err := query.SaveLedgerMember(ctx, sqlc.SaveLedgerMemberParams{
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

func mapLedgerError(err error) error {
	switch {
	case err == nil:
		return nil
	case isViolatingConstraint(err, constraintLedgerUniqueMember):
		return v1.ErrConflict
	case isViolatingConstraint(err, constraintLedgerMembersFK):
		return v1.ErrNotFound
	default:
		return mapError(err)
	}
}
