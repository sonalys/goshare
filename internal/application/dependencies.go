package application

import (
	"context"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	LedgerQueries interface {
		// Get returns the ledger by id.
		// Returns v1.ErrNotFound if it doesn't exist.
		Get(ctx context.Context, id domain.ID) (*domain.Ledger, error)
		// ListByUser returns all ledgers that the identity created or is a member.
		// Returns empty list if nothing is found.
		ListByUser(ctx context.Context, identity domain.ID) ([]domain.Ledger, error)
	}

	LedgerCommands interface {
		Create(ctx context.Context, ledger *domain.Ledger) error
		Update(ctx context.Context, ledger *domain.Ledger) error
	}

	LedgerRepository interface {
		LedgerQueries
		LedgerCommands
	}

	UserQueries interface {
		Get(ctx context.Context, id domain.ID) (*domain.User, error)
		GetByEmail(ctx context.Context, email string) (*domain.User, error)
		ListByEmail(ctx context.Context, emails []string) ([]domain.User, error)
	}

	UserCommands interface {
		Save(ctx context.Context, user *domain.User) error
	}

	UserRepository interface {
		UserQueries
		UserCommands
	}

	ExpenseQueries interface {
		Get(ctx context.Context, id domain.ID) (*domain.Expense, error)
		ListByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error)
	}

	ExpenseCommands interface {
		Create(ctx context.Context, ledgerID domain.ID, expense *domain.Expense) error
		Update(ctx context.Context, expense *domain.Expense) error
	}

	ExpenseRepository interface {
		ExpenseQueries
		ExpenseCommands
	}

	Database interface {
		Queries
		Transaction(ctx context.Context, f func(tx Repositories) error) error
	}

	Queries interface {
		Expense() ExpenseQueries
		Ledger() LedgerQueries
		User() UserQueries
	}

	Repositories interface {
		Expense() ExpenseRepository
		Ledger() LedgerRepository
		User() UserRepository
	}
)
