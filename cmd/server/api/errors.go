package api

import (
	"context"
	"net/http"

	"github.com/sonalys/goshare/cmd/server/handlers"
)

func newRespUnauthorized(ctx context.Context) *handlers.ErrorResponseStatusCode {
	return newErrorResponse(ctx, http.StatusUnauthorized, handlers.Error{
		Code:    handlers.ErrorCodeUnauthorized,
		Message: "unauthorized",
	})
}
