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
	return mapError(err)
}
