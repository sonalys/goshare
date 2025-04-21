package ledgers

import (
	"context"
	"time"

	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	LedgerRepository interface {
		Create(ctx context.Context, ledger *v1.Ledger) error
		GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error)
		GetLedgerBalance(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipantBalance, error)
		AddParticipant(ctx context.Context, ledgerID, userID, invitedUserID v1.ID) error
		AddParticipants(ctx context.Context, ledgerID, userID v1.ID, ids ...v1.ID) error
	}

	UserRepository interface {
		GetByEmail(ctx context.Context, emails []string) ([]v1.User, error)
	}

	ExpenseRepository interface {
		Create(ctx context.Context, expense *v1.Expense) error
		GetByLedger(ctx context.Context, ledgerID v1.ID, cursor time.Time, limit int32) ([]v1.Expense, error)
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
