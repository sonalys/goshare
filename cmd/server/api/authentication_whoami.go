package api

import (
	"context"
	"log/slog"

	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
)

// GetIdentity implements handlers.StrictServerInterface.
func (a *API) GetIdentity(ctx context.Context, request handlers.GetIdentityRequestObject) (handlers.GetIdentityResponseObject, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "could not retrieve identity", slog.Any("error", err))
		return nil, err
	}

	return handlers.GetIdentity200JSONResponse{
		Email:  types.Email(identity.Email),
		UserId: identity.UserID,
	}, nil
}
