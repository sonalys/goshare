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
func (a *API) Login(ctx context.Context, req *handlers.LoginReq) (*handlers.LoginOK, error) {
	resp, err := a.dependencies.UserController.Login(ctx, users.LoginRequest{
		Email:    string(req.Email),
		Password: req.Password,
	})
	switch {
	case err == nil:
		return &handlers.LoginOK{
			SetCookie: handlers.NewOptString(fmt.Sprintf("SESSIONID=%s; Path=/; HttpOnly; SameSite=Strict", resp.Token)),
		}, nil
	case errors.Is(err, v1.ErrEmailPasswordMismatch):
		return nil, newErrorResponse(ctx, http.StatusUnauthorized, handlers.Error{
			Code:    handlers.ErrorCodeEmailPasswordMismatch,
			Message: "invalid credentials",
		})
	default:
		return nil, err
	}
}
