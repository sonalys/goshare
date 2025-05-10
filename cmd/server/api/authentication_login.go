package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

// Login implements handlers.StrictServerInterface.
func (a *API) AuthenticationLogin(ctx context.Context, req *handlers.AuthenticationLoginReq) (*handlers.AuthenticationLoginOK, error) {
	resp, err := a.Users.Login(ctx, controllers.LoginRequest{
		Email:    string(req.Email),
		Password: req.Password,
	})
	if err == nil {
		return &handlers.AuthenticationLoginOK{
			SetCookie: handlers.NewOptString(fmt.Sprintf("SESSIONID=%s; Path=/; HttpOnly; SameSite=Strict", resp.Token)),
		}, nil
	}

	if target := new(v1.ErrUserCredentialsMismatch); errors.As(err, &target) {
		return nil, newErrorResponse(ctx, http.StatusUnauthorized, handlers.Error{
			Code:    handlers.ErrorCodeEmailPasswordMismatch,
			Message: "invalid credentials",
		})
	}

	return nil, err
}
