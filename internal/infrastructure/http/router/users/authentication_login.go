package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sonalys/goshare/internal/application/controllers/identitycontroller"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	v1 "github.com/sonalys/goshare/pkg/v1"
)

func (a *Router) AuthenticationLogin(ctx context.Context, req *server.AuthenticationLoginReq) (*server.AuthenticationLoginOK, error) {
	resp, err := a.IdentityController.Login(ctx, identitycontroller.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == nil {
		return &server.AuthenticationLoginOK{
			SetCookie: server.NewOptString(fmt.Sprintf("SESSIONID=%s; Path=/; HttpOnly; SameSite=Strict", resp.Token)),
		}, nil
	}

	if target := new(v1.UserCredentialsMismatchError); errors.As(err, &target) {
		return nil, &server.ErrorResponseStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response: server.ErrorResponse{
				Errors: []server.Error{
					server.Error{
						Code:    server.ErrorCodeEmailPasswordMismatch,
						Message: "username or password invalid",
						Metadata: server.NewOptErrorMetadata(server.ErrorMetadata{
							Field: server.NewOptString("email"),
						}),
					},
				},
			},
		}
	}

	return nil, err
}
