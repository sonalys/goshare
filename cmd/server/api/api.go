package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/cmd/server/handlers"
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

	UserLedgerLister interface {
		ListByUser(ctx context.Context, userID uuid.UUID) ([]v1.Ledger, error)
	}

	Dependencies struct {
		UserRegister
		UserAuthentication
		LedgerCreater
		UserLedgerLister
	}

	API struct {
		dependencies Dependencies
	}
)

var (
	_ handlers.StrictServerInterface = (*API)(nil)
)

func New(dependencies Dependencies) *API {
	return &API{
		dependencies: dependencies,
	}
}
