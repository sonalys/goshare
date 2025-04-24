package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type LedgerRepository struct {
	client *Client
}

func NewLedgerRepository(client *Client) *LedgerRepository {
	return &LedgerRepository{
		client: client,
	}
}

func (r *LedgerRepository) Create(ctx context.Context, userID v1.ID, createFn func(count int64) (*v1.Ledger, error)) error {
	return mapLedgerError(r.client.transaction(ctx, func(tx pgx.Tx) error {
		query := r.client.queries().WithTx(tx)

		id := convertUUID(userID)

		if err := query.LockUserForUpdate(ctx, id); err != nil {
			return fmt.Errorf("failed to acquire user lock for updating ledger")
		}

		userLedgersCount, err := query.CountUserLedgers(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to count user ledgers")
		}

		ledger, err := createFn(userLedgersCount)
		if err != nil {
			return fmt.Errorf("failed to create ledger: %w", err)
		}

		createLedgerReq := sqlc.CreateLedgerParams{
			ID:        convertUUID(ledger.ID),
			Name:      ledger.Name,
			CreatedAt: convertTime(ledger.CreatedAt),
			CreatedBy: convertUUID(ledger.CreatedBy),
		}

		if err := query.CreateLedger(ctx, createLedgerReq); err != nil {
			return fmt.Errorf("failed to create ledger: %w", err)
		}

		for _, participant := range ledger.Participants {
			addReq := sqlc.AddUserToLedgerParams{
				ID:        convertUUID(participant.ID),
				LedgerID:  createLedgerReq.ID,
				UserID:    convertUUID(participant.UserID),
				CreatedAt: createLedgerReq.CreatedAt,
				CreatedBy: createLedgerReq.CreatedBy,
			}

			if err := query.AddUserToLedger(ctx, addReq); err != nil {
				return fmt.Errorf("failed to add user to ledger: %w", err)
			}
		}

		return nil
	}))
}

func (r *LedgerRepository) Find(ctx context.Context, id v1.ID) (*v1.Ledger, error) {
	ledger, err := r.client.queries().FindLedgerById(ctx, convertUUID(id))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	participants, err := r.client.queries().GetLedgerParticipants(ctx, convertUUID(id))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	return mappers.NewLedger(&ledger, participants), nil
}

func (r *LedgerRepository) GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error) {
	ledgers, err := r.client.queries().GetUserLedgers(ctx, convertUUID(userID))
	if err != nil {
		return nil, mapLedgerError(err)
	}

	result := make([]v1.Ledger, 0, len(ledgers))
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
		return v1.ErrUserAlreadyMember
	case isViolatingConstraint(err, constraintLedgerParticipantsFK):
		return v1.ErrNotFound
	default:
		return mapError(err)
	}
}
