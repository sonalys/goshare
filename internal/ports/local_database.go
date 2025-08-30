package ports

import (
	"context"
)

type (
	Repositories interface {
		Expense() ExpenseRepository
		Ledger() LedgerRepository
		User() UserRepository
	}

	LocalDatabase interface {
		Expense() ExpenseQueries
		Ledger() LedgerQueries
		User() UserQueries
		Transaction(ctx context.Context, f func(tx Repositories) error) error
	}
)
