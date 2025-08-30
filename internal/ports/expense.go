package ports

import (
	"context"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
)

type (
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
)
