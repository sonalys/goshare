package repositories

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
)

var ledgerConstraintMapping = map[string]error{
	"fk_ledger_created_by":        domain.ErrUserNotFound,
	"fk_ledger_member_ledger":     domain.ErrLedgerNotFound,
	"fk_ledger_member_user":       domain.ErrUserNotFound,
	"fk_ledger_member_created_by": domain.ErrUserNotFound,
	"unique_ledger_member":        fmt.Errorf("member already in ledger: %w", domain.ErrConflict),
}

type LedgerRepository struct {
	client postgres.Connection
}

func newLedgerRepository(client postgres.Connection) *LedgerRepository {
	return &LedgerRepository{
		client: client,
	}
}

func (r *LedgerRepository) transaction(ctx context.Context, f func(q postgres.Connection) error) error {
	return ledgerError(r.client.Transaction(ctx, f))
}

func ledgerError(err error) error {
	if err := postgres.MapConstraintError(err, ledgerConstraintMapping); err != nil {
		return err
	}

	return postgres.DefaultErrorMapping(err)
}
