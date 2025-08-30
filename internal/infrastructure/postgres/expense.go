package postgres

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/domain"
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
	client connection
}

func (r *ExpenseRepository) transaction(ctx context.Context, f func(q connection) error) error {
	return expenseError(r.client.transaction(ctx, f))
}

func expenseError(err error) error {
	if err := constraintErrorMap(err, expenseConstraintMapping); err != nil {
		return err
	}

	return defaultErrorMapping(err)
}
