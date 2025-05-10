package controllers

import (
	"context"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	LedgerRepository interface {
		Create(ctx context.Context, ledger *domain.Ledger) error
		Find(ctx context.Context, id domain.ID) (*domain.Ledger, error)
		GetByUser(ctx context.Context, identity domain.ID) ([]domain.Ledger, error)
		Update(ctx context.Context, ledger *domain.Ledger) error
	}

	UserRepository interface {
		Create(ctx context.Context, user *domain.User) error
		Find(ctx context.Context, id domain.ID) (*domain.User, error)
		FindByEmail(ctx context.Context, email string) (*domain.User, error)
		ListByEmail(ctx context.Context, emails []string) ([]domain.User, error)
	}

	ExpenseRepository interface {
		Create(ctx context.Context, ledgerID domain.ID, expense *domain.Expense) error
		Find(ctx context.Context, id domain.ID) (*domain.Expense, error)
		GetByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error)
		Update(ctx context.Context, expense *domain.Expense) error
	}

	Database interface {
		Repositories
		Transaction(ctx context.Context, f func(db Database) error) error
	}

	Repositories interface {
		Expense() ExpenseRepository
		Ledger() LedgerRepository
		User() UserRepository
	}

	IdentityEncoder interface {
		Encode(identity *v1.Identity) (string, error)
	}

	Dependencies struct {
		Database
		IdentityEncoder
	}
)
