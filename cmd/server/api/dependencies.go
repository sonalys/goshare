package api

import (
	"context"

	"github.com/sonalys/goshare/internal/application/ledgers"
	"github.com/sonalys/goshare/internal/application/users"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	UserController interface {
		Register(ctx context.Context, req users.RegisterRequest) (*users.RegisterResponse, error)
		Login(ctx context.Context, req users.LoginRequest) (*users.LoginResponse, error)
	}

	LedgerController interface {
		Create(ctx context.Context, req ledgers.CreateRequest) (*ledgers.CreateResponse, error)
		GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error)
		AddParticipants(ctx context.Context, req ledgers.AddMembersRequest) error
		GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error)
		CreateExpense(ctx context.Context, req ledgers.CreateExpenseRequest) (*ledgers.CreateExpenseResponse, error)
	}

	Dependencies struct {
		LedgerController
		UserController
	}
)
