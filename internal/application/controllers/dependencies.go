package controllers

import (
	"context"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

type (
	LedgerRepository interface {
		Create(ctx context.Context, userID v1.ID, createFn func(count int64) (*v1.Ledger, error)) error
		GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error)
		AddParticipants(ctx context.Context, ledgerID v1.ID, updateFn func(*v1.Ledger) error) error
		GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error)
		Find(ctx context.Context, id v1.ID) (*v1.Ledger, error)
	}

	UserRepository interface {
		ListByEmail(ctx context.Context, emails []string) ([]v1.User, error)
		Create(ctx context.Context, user *v1.User) error
		FindByEmail(ctx context.Context, email string) (*v1.User, error)
	}

	ExpenseRepository interface {
		Create(ctx context.Context, ledgerID v1.ID, createFn func(ledger *v1.Ledger) (*v1.Expense, error)) error
		Find(ctx context.Context, id v1.ID) (*v1.Expense, error)
		GetByLedger(ctx context.Context, ledgerID v1.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error)
	}

	Database interface {
		Repositories
		Transaction(ctx context.Context, f func(r Repositories) error) error
	}

	Repositories interface {
		Ledger() LedgerRepository
		User() UserRepository
		Expense() ExpenseRepository
	}

	IdentityEncoder interface {
		Encode(identity *v1.Identity) (string, error)
	}

	Dependencies struct {
		Database
		IdentityEncoder
	}
)
