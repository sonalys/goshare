package repositories

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres"
)

var expenseConstraintMapping = map[string]error{
	"fk_expense_ledger":            domain.ErrLedgerNotFound,
	"fk_expense_created_by":        domain.ErrUserNotFound,
	"fk_expense_updated_by":        domain.ErrUserNotFound,
	"fk_expense_record_expense":    domain.ErrExpenseNotFound,
	"fk_expense_record_from_user":  domain.ErrUserNotFound,
	"fk_expense_record_to_user":    domain.ErrUserNotFound,
	"fk_expense_record_created_by": domain.ErrUserNotFound,
	"unique_expense_record":        fmt.Errorf("expense already exists: %w", domain.ErrConflict),
}

type ExpenseRepository struct {
	conn postgres.Connection
}

func newExpenseRepository(conn postgres.Connection) *ExpenseRepository {
	return &ExpenseRepository{
		conn: conn,
	}
}

func (r *ExpenseRepository) transaction(ctx context.Context, f func(q postgres.Connection) error) error {
	return expenseError(r.conn.Transaction(ctx, f))
}

func expenseError(err error) error {
	if err := postgres.MapConstraintError(err, expenseConstraintMapping); err != nil {
		return err
	}

	return postgres.DefaultErrorMapping(err)
}
