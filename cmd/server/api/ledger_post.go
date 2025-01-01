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
	identity, err := GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	req := ledgers.CreateRequest{
		UserID: identity.UserID,
		Name:   request.Body.Name,
	}

	switch resp, err := a.dependencies.LedgerCreater.CreateLedger(ctx, req); {
	case err != nil:
		return handlers.CreateLedger200JSONResponse{Id: resp.ID}, nil
	default:
		if errList := new(v1.FieldErrorList); errors.As(err, errList) {
			return handlers.CreateLedgerdefaultJSONResponse{
				Body:       newErrorResponse(ctx, getCausesFromFieldErrors(*errList)),
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		return nil, err
	}
}
