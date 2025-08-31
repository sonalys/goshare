package ledger

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
)

var constraintMapping = map[string]error{
	"fk_ledger_created_by":        domain.ErrUserNotFound,
	"fk_ledger_member_ledger":     domain.ErrLedgerNotFound,
	"fk_ledger_member_user":       domain.ErrUserNotFound,
	"fk_ledger_member_created_by": domain.ErrUserNotFound,
	"unique_ledger_member":        fmt.Errorf("member already in ledger: %w", domain.ErrConflict),
}

type Repository struct {
	client postgres.Connection
}

func New(client postgres.Connection) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) transaction(ctx context.Context, f func(q postgres.Connection) error) error {
	return ledgerError(r.client.Transaction(ctx, f))
}

func ledgerError(err error) error {
	if err == nil {
		return nil
	}

	if err := postgres.MapConstraintError(err, constraintMapping); err != nil {
		return err
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrLedgerNotFound
	default:
		return postgres.DefaultErrorMapping(err)
	}
}
