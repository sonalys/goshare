package ledgers

import (
	"context"

	"github.com/google/uuid"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	LedgerRepository interface {
		Create(ctx context.Context, ledger *v1.Ledger) error
		GetByUser(ctx context.Context, userID uuid.UUID) ([]v1.Ledger, error)
		GetLedgerBalance(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerParticipantBalance, error)
	}

	ExpenseRepository interface {
		Create(ctx context.Context, expense *v1.Expense) error
	}

	Controller struct {
		ledgerRepository  LedgerRepository
		expenseRepository ExpenseRepository
	}
)

func NewController(
	ledgerRepository LedgerRepository,
	expenseRepository ExpenseRepository,
) *Controller {
	return &Controller{
		ledgerRepository:  ledgerRepository,
		expenseRepository: expenseRepository,
	}
}
