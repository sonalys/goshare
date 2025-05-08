package postgres

import (
	"context"
	"fmt"

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

func (r *LedgerRepository) transaction(ctx context.Context, f func(q *sqlc.Queries) error) error {
	return mapLedgerError(r.client.transaction(ctx, f))
}

func (r *LedgerRepository) Create(ctx context.Context, ledger *domain.Ledger) error {
	return r.transaction(ctx, func(query *sqlc.Queries) error {
		createLedgerReq := sqlc.CreateLedgerParams{
			ID:        convertID(ledger.ID),
			Name:      ledger.Name,
			CreatedAt: convertTime(ledger.CreatedAt),
			CreatedBy: convertID(ledger.CreatedBy),
		}

		if err := query.CreateLedger(ctx, createLedgerReq); err != nil {
			return fmt.Errorf("failed to create ledger: %w", err)
		}

		for _, participant := range ledger.Participants {
			addReq := sqlc.AddUserToLedgerParams{
				ID:        convertID(participant.ID),
				LedgerID:  createLedgerReq.ID,
				UserID:    convertID(participant.Identity),
				CreatedAt: createLedgerReq.CreatedAt,
				CreatedBy: createLedgerReq.CreatedBy,
			}

			if err := query.AddUserToLedger(ctx, addReq); err != nil {
				return fmt.Errorf("failed to add user to ledger: %w", err)
			}
		}

		return nil
	})
}

func (r *LedgerRepository) Find(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
	ledger, err := r.client.queries().FindLedgerById(ctx, convertID(id))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	participants, err := r.client.queries().GetLedgerParticipants(ctx, convertID(id))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	return mappers.NewLedger(&ledger, participants), nil
}

func (r *LedgerRepository) GetByUser(ctx context.Context, userID domain.ID) ([]domain.Ledger, error) {
	ledgers, err := r.client.queries().GetUserLedgers(ctx, convertID(userID))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	result := make([]domain.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		participants, err := r.client.queries().GetLedgerParticipants(ctx, ledger.ID)
		if err != nil {
			return nil, mapLedgerError(err)
		}
		result = append(result, *mappers.NewLedger(&ledger, participants))
	}
	return result, nil
}

func mapLedgerError(err error) error {
	switch {
	case err == nil:
		return nil
	case isViolatingConstraint(err, constraintLedgerUniqueParticipant):
		return domain.ErrUserAlreadyMember
	case isViolatingConstraint(err, constraintLedgerParticipantsFK):
		return domain.ErrNotFound
	default:
		return mapError(err)
	}
}
