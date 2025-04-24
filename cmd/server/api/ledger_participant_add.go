package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) LedgerParticipantAdd(ctx context.Context, req *handlers.LedgerParticipantAddReq, params handlers.LedgerParticipantAddParams) error {
	identity, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	apiParams := ledgers.AddMembersRequest{
		UserID:   identity.UserID,
		LedgerID: v1.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	}

	return a.dependencies.LedgerController.AddParticipants(ctx, apiParams)
}
