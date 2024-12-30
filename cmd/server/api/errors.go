package api

import (
	"github.com/sonalys/goshare/cmd/server/handlers"
)

func newErrorResponse(cause []handlers.Error) handlers.ErrorResponseJSONResponse {
	return handlers.ErrorResponseJSONResponse{
		Errors: cause,
	}
}
