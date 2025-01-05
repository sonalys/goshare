package api

import (
	"context"
	"errors"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
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
		LedgerID: v1.ConvertID(request.LedgerID),
		Emails:   collectEmails(request.Body.Emails),
	}
	switch err := a.dependencies.LedgerMemberCreater.AddMembers(ctx, req); {
	case err == nil:
		return handlers.AddLedgerMember202Response{}, nil
	case errors.Is(err, v1.ErrLedgerMaxUsers):
		return handlers.AddLedgerMemberdefaultJSONResponse{
			Body: newErrorResponse(ctx, []handlers.Error{
				{
					Code:    handlers.LedgerMaxUsers,
					Message: err.Error(),
				},
			}),
			StatusCode: http.StatusBadRequest,
		}, nil
	default:
		if errList := new(v1.FieldErrorList); errors.As(err, errList) {
			return handlers.AddLedgerMemberdefaultJSONResponse{
				Body:       newErrorResponse(ctx, getCausesFromFieldErrors(*errList)),
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		return nil, err
	}
}
