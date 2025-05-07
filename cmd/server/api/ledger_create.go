package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
)

func (a *API) LedgerCreate(ctx context.Context, req *handlers.LedgerCreateReq) (r *handlers.LedgerCreateOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	apiParams := controllers.CreateLedgerRequest{
		Identity: identity.UserID,
		Name:     req.Name,
	}

	switch resp, err := a.Ledgers.Create(ctx, apiParams); err {
	case nil:
		return &handlers.LedgerCreateOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}
