package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/ledgers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (a *API) CreateLedger(ctx context.Context, request handlers.CreateLedgerRequestObject) (handlers.CreateLedgerResponseObject, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	req := ledgers.CreateRequest{
		UserID: identity.UserID,
		Name:   request.Body.Name,
	}

	switch resp, err := a.dependencies.LedgerCreater.Create(ctx, req); {
	case err == nil:
		return handlers.CreateLedger200JSONResponse{Id: resp.ID.UUID()}, nil
	case errors.Is(err, v1.ErrUserMaxLedgers):
		return handlers.CreateLedgerdefaultJSONResponse{
			Body: newErrorResponse(ctx, []handlers.Error{
				{
					Code:    handlers.UserMaxLedgers,
					Message: err.Error(),
				},
			}),
			StatusCode: http.StatusBadRequest,
		}, nil
	default:
		if causes, ok := extractErrorCauses(err); ok {
			return handlers.CreateLedgerdefaultJSONResponse{
				Body:       newErrorResponse(ctx, getCausesFromFieldErrors(causes)),
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		return nil, err
	}
}
