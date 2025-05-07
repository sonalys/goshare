package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

// Login implements handlers.StrictServerInterface.
func (a *API) AuthenticationLogin(ctx context.Context, req *handlers.AuthenticationLoginReq) (*handlers.AuthenticationLoginOK, error) {
	resp, err := a.Users.Login(ctx, controllers.LoginRequest{
		Email:    string(req.Email),
		Password: req.Password,
	})
	switch {
	case err == nil:
		return &handlers.AuthenticationLoginOK{
			SetCookie: handlers.NewOptString(fmt.Sprintf("SESSIONID=%s; Path=/; HttpOnly; SameSite=Strict", resp.Token)),
		}, nil
	case errors.Is(err, domain.ErrEmailPasswordMismatch):
		return nil, newErrorResponse(ctx, http.StatusUnauthorized, handlers.Error{
			Code:    handlers.ErrorCodeEmailPasswordMismatch,
			Message: "invalid credentials",
		})
	default:
		return nil, err
	}
}
