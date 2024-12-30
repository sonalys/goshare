package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
)

type (
	Dependencies struct{}

	API struct {
		dependencies Dependencies
	}
)

var (
	_ handlers.StrictServerInterface = (*API)(nil)
)

func New(dependencies Dependencies) *API {
	return &API{
		dependencies: dependencies,
	}
}

func (a *API) GetHealthcheck(ctx context.Context, _ handlers.GetHealthcheckRequestObject) (handlers.GetHealthcheckResponseObject, error) {
	return handlers.GetHealthcheck200JSONResponse{}, nil
}
