package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerMemberAdd(ctx context.Context, req *handlers.LedgerMemberAddReq, params handlers.LedgerMemberAddParams) error {
	identity, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	apiParams := usercontroller.AddMembersRequest{
		Actor:    identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	}

	return a.UserController.AddMembers(ctx, apiParams)
}
