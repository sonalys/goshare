package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerCreate(ctx context.Context, req *server.LedgerCreateReq) (r *server.LedgerCreateOK, _ error) {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	apiParams := usercontroller.CreateLedgerRequest{
		ActorID: identity.UserID,
		Name:    req.Name,
	}

	resp, err := a.Ledgers().Create(ctx, apiParams)
	if err != nil {
		return nil, err
	}

	return &server.LedgerCreateOK{
		ID: resp.ID.UUID(),
	}, nil
}
