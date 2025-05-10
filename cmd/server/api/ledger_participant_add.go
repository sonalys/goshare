package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerMemberAdd(ctx context.Context, req *handlers.LedgerMemberAddReq, params handlers.LedgerMemberAddParams) error {
	identity, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	apiParams := controllers.AddMembersRequest{
		Actor:    identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	}

	return a.Ledgers.AddMembers(ctx, apiParams)
}
