package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

func (a *API) UserRegister(ctx context.Context, req *handlers.UserRegisterReq) (r *handlers.UserRegisterOK, _ error) {
	apiParams := identitycontroller.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     string(req.Email),
		Password:  req.Password,
	}

	switch resp, err := a.IdentityController.Register(ctx, apiParams); {
	case err == nil:
		return &handlers.UserRegisterOK{
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
