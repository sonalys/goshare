package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerCreate(ctx context.Context, req *server.LedgerCreateReq) (r *server.LedgerCreateOK, _ error) {
	identity, err := middlewares.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	apiParams := usercontroller.CreateLedgerRequest{
		ActorID: identity.UserID,
		Name:    req.Name,
	}

	switch resp, err := a.UserController.Ledgers().Create(ctx, apiParams); err {
	case nil:
		return &server.LedgerCreateOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}
