package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
)

func (a *API) CreateLedger(ctx context.Context, req *handlers.CreateLedgerReq) (r *handlers.CreateLedgerOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	apiParams := ledgers.CreateRequest{
		UserID: identity.UserID,
		Name:   req.Name,
	}

	switch resp, err := a.dependencies.LedgerCreater.Create(ctx, apiParams); {
	case err == nil:
		return &handlers.CreateLedgerOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}
