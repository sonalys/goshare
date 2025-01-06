package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/users"
	"github.com/sonalys/goshare/internal/pkg/pointers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) RegisterUser(ctx context.Context, request handlers.RegisterUserRequestObject) (handlers.RegisterUserResponseObject, error) {
	req := users.RegisterRequest{
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
		Email:     string(request.Body.Email),
		Password:  request.Body.Password,
	}

	switch resp, err := a.dependencies.UserRegister.Register(ctx, req); {
	case err == nil:
		return handlers.RegisterUser200JSONResponse{Id: resp.ID.UUID()}, nil
	case errors.Is(err, v1.ErrEmailAlreadyRegistered):
		return handlers.RegisterUserdefaultJSONResponse{
			Body: newErrorResponse(ctx, []handlers.Error{
				{
					Code:    handlers.InvalidField,
					Message: err.Error(),
					Metadata: &handlers.ErrorMetadata{
						Field: pointers.New("email"),
					},
				},
			}),
			StatusCode: http.StatusConflict,
		}, nil
	default:
		if errList := new(v1.FieldErrorList); errors.As(err, errList) {
			return handlers.RegisterUserdefaultJSONResponse{
				Body:       newErrorResponse(ctx, getCausesFromFieldErrors(*errList)),
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		return nil, err
	}
}
