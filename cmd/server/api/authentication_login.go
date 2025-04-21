package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/users"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

// Login implements handlers.StrictServerInterface.
func (a *API) Login(ctx context.Context, req *handlers.LoginReq) (*handlers.LoginOK, error) {
	resp, err := a.dependencies.UserAuthentication.Login(ctx, users.LoginRequest{
		Email:    string(req.Email),
		Password: req.Password,
	})
	switch {
	case err == nil:
		return &handlers.LoginOK{
			SetCookie: handlers.NewOptString(fmt.Sprintf("SESSIONID=%s; Path=/; HttpOnly; SameSite=Strict", resp.Token)),
		}, nil
	case errors.Is(err, v1.ErrEmailPasswordMismatch):
		return nil, err
	default:
		return nil, err
	}
}
