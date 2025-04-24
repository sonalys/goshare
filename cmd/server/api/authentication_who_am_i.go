package api

import (
	"context"
	"log/slog"

	"github.com/sonalys/goshare/cmd/server/handlers"
)

// AuthenticationWhoAmI implements handlers.StrictServerInterface.
func (a *API) AuthenticationWhoAmI(ctx context.Context) (*handlers.AuthenticationWhoAmIOK, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "could not retrieve identity", slog.Any("error", err))
		return nil, err
	}

	return &handlers.AuthenticationWhoAmIOK{
		Email:  identity.Email,
		UserID: identity.UserID.UUID(),
	}, nil
}
