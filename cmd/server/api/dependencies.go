package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/application/ledgers"
	"github.com/sonalys/goshare/internal/application/users"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	UserRegister interface {
		Register(ctx context.Context, req users.RegisterRequest) (*users.RegisterResponse, error)
	}

	UserAuthentication interface {
		Login(ctx context.Context, req users.LoginRequest) (*users.LoginResponse, error)
	}

	LedgerCreater interface {
		Create(ctx context.Context, req ledgers.CreateRequest) (*ledgers.CreateResponse, error)
	}

	LedgerBalancesLister interface {
		GetBalances(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerParticipantBalance, error)
	}

	LedgerMemberCreater interface {
		AddMembers(ctx context.Context, req ledgers.AddMembersRequest) error
	}

	UserLedgerLister interface {
		GetByUser(ctx context.Context, userID uuid.UUID) ([]v1.Ledger, error)
	}

	ExpenseCreater interface {
		CreateExpense(ctx context.Context, req ledgers.CreateExpenseRequest) (*ledgers.CreateExpenseResponse, error)
	}

	Dependencies struct {
		ExpenseCreater
		LedgerBalancesLister
		LedgerCreater
		LedgerMemberCreater
		UserAuthentication
		UserLedgerLister
		UserRegister
	}
)
