package api

import (
	"context"
	"log/slog"

	"github.com/sonalys/goshare/cmd/server/handlers"
)

// GetIdentity implements handlers.StrictServerInterface.
func (a *API) GetIdentity(ctx context.Context) (*handlers.GetIdentityOK, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "could not retrieve identity", slog.Any("error", err))
		return nil, err
	}

	return &handlers.GetIdentityOK{
		Email:  identity.Email,
		UserID: identity.UserID.UUID(),
	}, nil
}
