package api

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
)

func collectEmails(from []openapi_types.Email) []string {
	to := make([]string, 0, len(from))

	for i := range from {
		to = append(to, string(from[i]))
	}

	return to
}

func (a *API) AddLedgerMember(ctx context.Context, request handlers.AddLedgerMemberRequestObject) (handlers.AddLedgerMemberResponseObject, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	req := ledgers.AddMembersRequest{
		UserID:   identity.UserID,
		LedgerID: request.LedgerID,
		Emails:   collectEmails(request.Body.Emails),
	}
	switch err := a.dependencies.LedgerMemberCreater.AddMembers(ctx, req); {
	case err == nil:
		return handlers.AddLedgerMember202Response{}, nil
	default:
		return nil, err
	}
}
