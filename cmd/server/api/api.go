package api

import (
	"context"
	"log/slog"

	"github.com/oapi-codegen/runtime/types"
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

// GetIdentity implements handlers.StrictServerInterface.
func (a *API) GetIdentity(ctx context.Context, request handlers.GetIdentityRequestObject) (handlers.GetIdentityResponseObject, error) {
	identity, err := GetIdentity(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "could not retrieve identity", slog.Any("error", err))
		return nil, err
	}

	return handlers.GetIdentity200JSONResponse{
		Email:  types.Email(identity.Email),
		UserId: identity.UserID,
	}, nil
}

var (
	_ handlers.StrictServerInterface = (*API)(nil)
)

func New(dependencies Dependencies) *API {
	return &API{
		dependencies: dependencies,
	}
}
