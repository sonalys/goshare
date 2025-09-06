package members

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerMemberAdd(ctx context.Context, req *server.LedgerMemberAddReq, params server.LedgerMemberAddParams) error {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return err
	}

	return a.Ledgers().MembersAdd(ctx, usercontroller.AddMembersRequest{
		ActorID:  identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	})
}
