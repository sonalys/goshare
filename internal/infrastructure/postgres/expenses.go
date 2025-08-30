package postgres

import (
	"context"
)

type ExpenseRepository struct {
	client connection
}

func (r *ExpenseRepository) transaction(ctx context.Context, f func(q connection) error) error {
	return mapExpensesError(r.client.transaction(ctx, f))
}

func mapExpensesError(err error) error {
	switch {
	case err == nil:
		return nil
	default:
		return mapError(err)
	}
}
