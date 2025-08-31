package expense

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
)

var constraintMapping = map[string]error{
	"fk_expense_ledger":            domain.ErrLedgerNotFound,
	"fk_expense_created_by":        domain.ErrUserNotFound,
	"fk_expense_updated_by":        domain.ErrUserNotFound,
	"fk_expense_record_expense":    domain.ErrExpenseNotFound,
	"fk_expense_record_from_user":  domain.ErrUserNotFound,
	"fk_expense_record_to_user":    domain.ErrUserNotFound,
	"fk_expense_record_created_by": domain.ErrUserNotFound,
	"unique_expense_record":        fmt.Errorf("expense already exists: %w", domain.ErrConflict),
}

type Repository struct {
	conn postgres.Connection
}

func New(conn postgres.Connection) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) transaction(ctx context.Context, f func(q postgres.Connection) error) error {
	return expenseError(r.conn.Transaction(ctx, f))
}

func expenseError(err error) error {
	if err == nil {
		return nil
	}

	if err := postgres.MapConstraintError(err, constraintMapping); err != nil {
		return err
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrExpenseNotFound
	default:
		return postgres.DefaultErrorMapping(err)
	}
}
