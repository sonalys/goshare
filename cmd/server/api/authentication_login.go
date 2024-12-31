package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/users"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

// Login implements handlers.StrictServerInterface.
func (a *API) Login(ctx context.Context, req handlers.LoginRequestObject) (handlers.LoginResponseObject, error) {
	resp, err := a.dependencies.UserAuthentication.Login(ctx, users.LoginRequest{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
	})
	switch {
	case err == nil:
		return handlers.Login200Response{
			Headers: handlers.Login200ResponseHeaders{
				SetCookie: fmt.Sprintf("SESSIONID=%s; Path=/; HttpOnly; SameSite=Strict", resp.Token),
			},
		}, nil
	case errors.Is(err, v1.ErrEmailPasswordMismatch):
		return handlers.LogindefaultJSONResponse{
			Body: newErrorResponse(ctx, []handlers.Error{
				{
					Code:    handlers.EmailPasswordMismatch,
					Message: err.Error(),
				},
			}),
			StatusCode: http.StatusForbidden,
		}, nil
	default:
		if errList := new(v1.FieldErrorList); errors.As(err, errList) {
			return handlers.LogindefaultJSONResponse{
				Body:       newErrorResponse(ctx, getCausesFromFieldErrors(*errList)),
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		return nil, err
	}

}
