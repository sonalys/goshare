package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

func (a *API) LedgerParticipantAdd(ctx context.Context, req *handlers.LedgerParticipantAddReq, params handlers.LedgerParticipantAddParams) error {
	identity, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	apiParams := controllers.AddMembersRequest{
		UserID:   identity.UserID,
		LedgerID: v1.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	}

	return a.Ledgers.AddParticipants(ctx, apiParams)
}
