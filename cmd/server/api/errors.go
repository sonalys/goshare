package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/pkg/pointers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/trace"
)

const urlKey = "url-key"

func InjectRequestContextDataMiddleware(handler nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		ctx = context.WithValue(ctx, urlKey, r.URL.Path)
		return handler(ctx, w, r, request)
	}
}

func getURL(ctx context.Context) string {
	url, _ := ctx.Value(urlKey).(string)
	return url
}

func newErrorResponse(ctx context.Context, cause []handlers.Error) handlers.ErrorResponseJSONResponse {
	return handlers.ErrorResponseJSONResponse{
		TraceId: types.UUID(trace.SpanContextFromContext(ctx).TraceID()),
		Url:     getURL(ctx),
		Errors:  cause,
	}
}

func getFieldErrorCode(from v1.FieldError) handlers.ErrorCode {
	switch {
	case errors.Is(from.Cause, v1.ErrInvalidValue):
		return handlers.InvalidField
	case errors.Is(from.Cause, v1.ErrRequiredValue):
		return handlers.RequiredField
	default:
		return handlers.ErrorCode("")
	}
}

func newCauseFromFieldError(from v1.FieldError) handlers.Error {
	return handlers.Error{
		Message: from.Error(),
		Code:    getFieldErrorCode(from),
		Metadata: &handlers.ErrorMetadata{
			Field: pointers.From(from.Field),
		},
	}
}

func getCausesFromFieldErrors(from v1.FieldErrorList) []handlers.Error {
	resp := make([]handlers.Error, 0, len(from))

	for i := range from {
		resp = append(resp, newCauseFromFieldError(from[i]))
	}

	return resp
}
