package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/users"
)

type (
	UserRegister interface {
		Register(ctx context.Context, req users.RegisterRequest) (*users.RegisterResponse, error)
	}

	UserAuthentication interface {
		Login(ctx context.Context, req users.LoginRequest) (*users.LoginResponse, error)
	}

	Dependencies struct {
		UserRegister
		UserAuthentication
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
