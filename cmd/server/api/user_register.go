package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/users"
)

func (a *API) RegisterUser(ctx context.Context, req *handlers.RegisterUserReq) (r *handlers.RegisterUserOK, _ error) {
	apiParams := users.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     string(req.Email),
		Password:  req.Password,
	}

	switch resp, err := a.dependencies.UserRegister.Register(ctx, apiParams); {
	case err == nil:
		return &handlers.RegisterUserOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}
