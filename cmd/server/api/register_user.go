package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/users"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) RegisterUser(ctx context.Context, req *handlers.RegisterUserReq) (r *handlers.RegisterUserOK, _ error) {
	apiParams := users.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     string(req.Email),
		Password:  req.Password,
	}

	switch resp, err := a.dependencies.UserController.Register(ctx, apiParams); {
	case err == nil:
		return &handlers.RegisterUserOK{
			ID: resp.ID.UUID(),
		}, nil
	case errors.Is(err, v1.ErrConflict):
		return nil, newErrorResponse(ctx, http.StatusConflict, handlers.Error{
			Code:    handlers.ErrorCodeInvalidField,
			Message: "already registered",
			Metadata: handlers.NewOptErrorMetadata(handlers.ErrorMetadata{
				Field: handlers.NewOptString("email"),
			}),
		})
	default:
		return nil, err
	}
}
