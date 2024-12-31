package api

import (
	"context"
	"errors"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/participants"
	"github.com/sonalys/goshare/internal/pkg/pointers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) RegisterUser(ctx context.Context, request handlers.RegisterUserRequestObject) (handlers.RegisterUserResponseObject, error) {
	resp, err := a.dependencies.ParticipantRegister.Register(ctx, participants.RegisterRequest{
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
		Email:     string(request.Body.Email),
		Password:  request.Body.Password,
	})
	switch {
	case err == nil:
		return handlers.RegisterUser200JSONResponse{Id: resp.ID}, nil
	case errors.Is(err, v1.ErrParticipantEmailAlreadyExists):
		return handlers.RegisterUser409JSONResponse{
			ErrorResponseJSONResponse: newErrorResponse([]handlers.Error{
				{
					Code:    handlers.InvalidField,
					Message: err.Error(),
					Metadata: &handlers.ErrorMetadata{
						Field: pointers.From("email"),
					},
				},
			}),
		}, nil
	default:
		return handlers.RegisterUserdefaultJSONResponse{}, err
	}
}
