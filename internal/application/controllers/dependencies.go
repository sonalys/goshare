package controllers

import (
	"context"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	LedgerRepository interface {
		Create(ctx context.Context, identity domain.ID, createFn func(count int64) (*domain.Ledger, error)) error
		GetByUser(ctx context.Context, identity domain.ID) ([]domain.Ledger, error)
		AddParticipants(ctx context.Context, ledgerID domain.ID, updateFn func(*domain.Ledger) error) error
		GetParticipants(ctx context.Context, ledgerID domain.ID) ([]domain.LedgerParticipant, error)
		Find(ctx context.Context, id domain.ID) (*domain.Ledger, error)
	}

	UserRepository interface {
		ListByEmail(ctx context.Context, emails []string) ([]domain.User, error)
		Create(ctx context.Context, user *domain.User) error
		FindByEmail(ctx context.Context, email string) (*domain.User, error)
	}

	ExpenseRepository interface {
		Create(ctx context.Context, ledgerID domain.ID, createFn func(ledger *domain.Ledger) (*domain.Expense, error)) error
		Find(ctx context.Context, id domain.ID) (*domain.Expense, error)
		GetByLedger(ctx context.Context, ledgerID domain.ID, cursor time.Time, limit int32) ([]v1.LedgerExpenseSummary, error)
	}

	Database interface {
		Repositories
		Transaction(ctx context.Context, f func(db Database) error) error
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
