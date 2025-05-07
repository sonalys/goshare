package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerParticipantAdd(ctx context.Context, req *handlers.LedgerParticipantAddReq, params handlers.LedgerParticipantAddParams) error {
	identity, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	apiParams := controllers.AddMembersRequest{
		Identity: identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	}

	return a.Ledgers.AddParticipants(ctx, apiParams)
}
