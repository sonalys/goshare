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
		AddParticipant(ctx context.Context, ledgerID, userID, invitedUserID uuid.UUID) error
	}

	UserRepository interface {
		GetByEmail(ctx context.Context, emails []string) ([]v1.User, error)
	}

	ExpenseRepository interface {
		Create(ctx context.Context, expense *v1.Expense) error
	}

	Controller struct {
		ledgerRepository  LedgerRepository
		expenseRepository ExpenseRepository
		userRepository    UserRepository
	}
)

func NewController(
	ledgerRepository LedgerRepository,
	expenseRepository ExpenseRepository,
	userReposiroty UserRepository,
) *Controller {
	return &Controller{
		ledgerRepository:  ledgerRepository,
		expenseRepository: expenseRepository,
		userRepository:    userReposiroty,
	}
}
