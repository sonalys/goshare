package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/queries"
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

func (r *LedgerRepository) Create(ctx context.Context, ledger *v1.Ledger) error {
	return mapLedgerError(r.client.queries().CreateLedger(ctx, queries.CreateLedgerParams{
		ID:        convertUUID(ledger.ID),
		Name:      ledger.Name,
		CreatedAt: convertTime(ledger.CreatedAt),
		CreatedBy: convertUUID(ledger.CreatedBy),
	}))
}

func (r *LedgerRepository) Find(ctx context.Context, id uuid.UUID) (*v1.Ledger, error) {
	ledger, err := r.client.queries().FindLedgerById(ctx, convertUUID(id))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	return newLedger(&ledger), nil
}

func (r *LedgerRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]v1.Ledger, error) {
	ledgers, err := r.client.queries().GetUserLedgers(ctx, convertUUID(userID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		result = append(result, *newLedger(&ledger))
	}
	return result, nil
}

func newLedger(ledger *queries.Ledger) *v1.Ledger {
	return &v1.Ledger{
		ID:        newUUID(ledger.ID),
		Name:      ledger.Name,
		CreatedAt: ledger.CreatedAt.Time,
		CreatedBy: newUUID(ledger.CreatedBy),
	}
}

func mapLedgerError(err error) error {
	switch err {
	case nil:
		return nil
	default:
		return mapError(err)
	}
}
