package users

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) AuthenticationWhoAmI(ctx context.Context) (*server.AuthenticationWhoAmIOK, error) {
	identity, err := middlewares.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return &server.AuthenticationWhoAmIOK{
		Email:  identity.Email,
		UserID: identity.UserID.UUID(),
	}, nil
}
