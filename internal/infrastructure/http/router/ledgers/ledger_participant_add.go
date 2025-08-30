package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/middlewares"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerMemberAdd(ctx context.Context, req *server.LedgerMemberAddReq, params server.LedgerMemberAddParams) error {
	identity, err := middlewares.GetIdentity(ctx)
	if err != nil {
		return err
	}

	apiParams := usercontroller.AddMembersRequest{
		ActorID:  identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	}

	return a.UserController.MembersAdd(ctx, apiParams)
}
