package users

import (
	"context"
	"errors"
	"net/http"

	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) UserRegister(ctx context.Context, req *server.UserRegisterReq) (r *server.UserRegisterOK, _ error) {
	apiParams := identitycontroller.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	switch resp, err := a.identityController.Register(ctx, apiParams); {
	case err == nil:
		return &server.UserRegisterOK{
			ID: resp.ID.UUID(),
		}, nil
	case errors.Is(err, domain.ErrUserAlreadyRegistered):
		return nil, &server.ErrorResponseStatusCode{
			StatusCode: http.StatusConflict,
			Response: server.ErrorResponse{
				Errors: []server.Error{
					{
						Code:    server.ErrorCodeInvalidField,
						Message: "already registered",
						Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
							Field: server.NewOptString("email"),
						}),
					},
				},
			},
		}
	default:
		return nil, err
	}
}
