package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/queries"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type ExpensesRepository struct {
	client *Client
}

func NewExpensesRepository(client *Client) *ExpensesRepository {
	return &ExpensesRepository{
		client: client,
	}
}

func (r *ExpensesRepository) Create(ctx context.Context, expense *v1.Expense) error {
	return mapError(r.client.transaction(ctx, func(tx *queries.Queries) error {
		return createExpense(ctx, tx, expense)
	}))
}

func (r *ExpensesRepository) Find(ctx context.Context, id uuid.UUID) (*v1.Expense, error) {
	expense, err := r.client.queries().FindExpenseById(ctx, convertUUID(id))
	if err != nil {
		return nil, mapError(err)
	}

	return newExpense(&expense), nil
}

func (r *ExpensesRepository) Update(ctx context.Context, expense *v1.Expense) error {
	return mapError(r.client.queries().UpdateExpense(ctx, queries.UpdateExpenseParams{
		ID:          convertUUID(expense.ID),
		Amount:      expense.Amount,
		CategoryID:  convertUUIDPtr(expense.CategoryID),
		Name:        expense.Name,
		ExpenseDate: convertTime(expense.ExpenseDate),
		UpdatedAt:   convertTime(expense.UpdatedAt),
		UpdatedBy:   convertUUID(expense.UpdatedBy),
	}))
}

func (r *ExpensesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return mapError(r.client.queries().DeleteExpense(ctx, convertUUID(id)))
}

func newExpense(expense *queries.Expense) *v1.Expense {
	return &v1.Expense{
		ID:          newUUID(expense.ID),
		Amount:      expense.Amount,
		CategoryID:  newUUIDPtr(expense.CategoryID),
		LedgerID:    newUUID(expense.LedgerID),
		Name:        expense.Name,
		ExpenseDate: expense.ExpenseDate.Time,
		CreatedAt:   expense.CreatedAt.Time,
		CreatedBy:   newUUID(expense.CreatedBy),
		UpdatedAt:   expense.UpdatedAt.Time,
		UpdatedBy:   newUUID(expense.UpdatedBy),
	}
}
