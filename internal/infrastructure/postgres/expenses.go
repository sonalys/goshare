package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	"github.com/sonalys/goshare/internal/pkg/monoids"
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
	return mapError(r.client.transaction(ctx, func(tx pgx.Tx) error {
		return r.createExpense(ctx, tx, expense)
	}))
}

func (r *ExpensesRepository) Update(ctx context.Context, expense *v1.Expense) error {
	return mapError(r.client.queries().UpdateExpense(ctx, sqlc.UpdateExpenseParams{
		ID:          convertUUID(expense.ID),
		Amount:      expense.Amount,
		CategoryID:  convertUUIDPtr(expense.CategoryID),
		Name:        expense.Name,
		ExpenseDate: convertTime(expense.ExpenseDate),
		UpdatedAt:   convertTime(expense.UpdatedAt),
		UpdatedBy:   convertUUID(expense.UpdatedBy),
	}))
}

func (r *ExpensesRepository) Delete(ctx context.Context, id v1.ID) error {
	return mapError(r.client.queries().DeleteExpense(ctx, convertUUID(id)))
}

func (r *ExpensesRepository) GetByLedger(ctx context.Context, ledgerID v1.ID, cursor time.Time, limit int32) ([]v1.Expense, error) {
	params := sqlc.GetLedgerExpensesParams{
		LedgerID:  convertUUID(ledgerID),
		CreatedAt: convertTime(cursor),
		Limit:     limit,
	}

	expenses, err := r.client.queries().GetLedgerExpenses(ctx, params)
	if err != nil {
		return nil, mapError(err)
	}

	ids := monoids.Map(expenses, func(from sqlc.Expense) pgtype.UUID {
		return from.ID
	})

	ledgerRecords, err := r.client.queries().GetExpensesRecords(ctx, ids)
	if err != nil {
		return nil, mapError(err)
	}

	return convertExpenses(expenses, ledgerRecords), nil
}

func newUserBalances(from []sqlc.LedgerRecord) []v1.ExpenseUserBalance {
	balances := make([]v1.ExpenseUserBalance, len(from))
	for i, record := range from {
		balances[i] = v1.ExpenseUserBalance{
			UserID:  newUUID(record.UserID),
			Balance: record.Amount,
		}
	}
	return balances
}

func newExpense(expense *sqlc.Expense, records []sqlc.LedgerRecord) *v1.Expense {
	return &v1.Expense{
		ID:           newUUID(expense.ID),
		Amount:       expense.Amount,
		CategoryID:   newUUIDPtr(expense.CategoryID),
		LedgerID:     newUUID(expense.LedgerID),
		Name:         expense.Name,
		ExpenseDate:  expense.ExpenseDate.Time,
		UserBalances: newUserBalances(records),
		CreatedAt:    expense.CreatedAt.Time,
		CreatedBy:    newUUID(expense.CreatedBy),
		UpdatedAt:    expense.UpdatedAt.Time,
		UpdatedBy:    newUUID(expense.UpdatedBy),
	}
}

func convertExpenses(from []sqlc.Expense, records []sqlc.LedgerRecord) []v1.Expense {
	to := make([]v1.Expense, 0, len(from))

	buffer := make([]sqlc.LedgerRecord, 0, len(records))
	for i := range from {
		buffer = buffer[:0]
		records := monoids.Reduce(records, buffer, func(acc []sqlc.LedgerRecord, record sqlc.LedgerRecord) []sqlc.LedgerRecord {
			if record.LedgerID == from[i].LedgerID {
				acc = append(acc, record)
			}
			return acc
		})

		to = append(to, *newExpense(&from[i], records))
	}

	return to
}
