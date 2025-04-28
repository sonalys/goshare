package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
)

// AuthenticationWhoAmI implements handlers.StrictServerInterface.
func (a *API) AuthenticationWhoAmI(ctx context.Context) (*handlers.AuthenticationWhoAmIOK, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return &handlers.AuthenticationWhoAmIOK{
		Email:  identity.Email,
		UserID: identity.UserID.UUID(),
	}, nil
}
