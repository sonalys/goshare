package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
)

func (a *API) GetHealthcheck(ctx context.Context, _ handlers.GetHealthcheckRequestObject) (handlers.GetHealthcheckResponseObject, error) {
	return handlers.GetHealthcheck200Response{}, nil
}
