package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) AddLedgerMember(ctx context.Context, req *handlers.AddLedgerMemberReq, params handlers.AddLedgerMemberParams) error {
	identity, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	apiParams := ledgers.AddMembersRequest{
		UserID:   identity.UserID,
		LedgerID: v1.ConvertID(params.LedgerID),
		Emails:   req.Emails,
	}
	switch err := a.dependencies.LedgerMemberCreater.AddMembers(ctx, apiParams); {
	case err == nil:
		return nil
	default:
		return err
	}
}
