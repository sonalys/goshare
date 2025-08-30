package postgres

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

var ledgerConstraintMapping = map[string]error{
	"fk_ledger_created_by":        domain.ErrUserNotFound,
	"fk_ledger_member_ledger":     domain.ErrLedgerNotFound,
	"fk_ledger_member_user":       domain.ErrUserNotFound,
	"fk_ledger_member_created_by": domain.ErrUserNotFound,
	"unique_ledger_member":        fmt.Errorf("member already in ledger: %w", domain.ErrConflict),
}

type LedgerRepository struct {
	client connection
}

func NewLedgerRepository(client connection) *LedgerRepository {
	return &LedgerRepository{
		client: client,
	}
}

func (r *LedgerRepository) transaction(ctx context.Context, f func(q *sqlcgen.Queries) error) error {
	return ledgerError(r.client.transaction(ctx, f))
}

func ledgerError(err error) error {
	if err := constraintErrorMap(err, ledgerConstraintMapping); err != nil {
		return err
	}

	return defaultErrorMapping(err)
}
