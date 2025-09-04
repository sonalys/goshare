package ports

import (
	"context"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	ExpenseQueries interface {
		// Get returns the expense by the given id.
		// Returns domain.ErrExpenseNotFound if it doesn't exist.
		Get(ctx context.Context, id domain.ID) (*domain.Expense, error)
		// ListByLedger returns all expenses from a ledger.
		// All results will be created before cursor, it's used for the pagination.
		// Limit restricts how many documents are returned.
		// If no results are found, an empty list and no error is returned.
		ListByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error)
	}

	ExpenseCommands interface {
		// Create will save the newly created expense.
		// Returns domain.ErrLedgerNotFound if the ledgerID doesn't exist.
		Create(ctx context.Context, expense *domain.Expense) error
		// Update will update the expense.
		// Returns domain.ErrLedgerNotFound if the expense doesn't exist.
		Update(ctx context.Context, expense *domain.Expense) error
	}

	ExpenseRepository interface {
		ExpenseQueries
		ExpenseCommands
	}
)
