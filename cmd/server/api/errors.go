package api

import (
	"context"
	"net/http"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/oapi-codegen/runtime/types"
	"github.com/sonalys/goshare/cmd/server/handlers"
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
