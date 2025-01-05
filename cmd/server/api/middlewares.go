package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

const urlKey = "url-key"

func InjectRequestContextDataMiddleware(handler nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
		ctx = context.WithValue(ctx, urlKey, r.URL.Path)
		return handler(ctx, w, r, request)
	}
}

func getURL(ctx context.Context) string {
	url, _ := ctx.Value(urlKey).(string)
	return url
}

const userCtxKey = "user-key"

func getIdentity(ctx context.Context) (*v1.Identity, error) {
	identity, ok := ctx.Value(userCtxKey).(*v1.Identity)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return identity, nil
}

var UnauthorizedResp = []handlers.Error{
	{
		Code:    handlers.Unauthorized,
		Message: "The request authentication failed",
	},
}

var AuthorizationExpiredResp = []handlers.Error{
	{
		Code:    handlers.Unauthorized,
		Message: "The request authentication is expired",
	},
}

var ForbiddenResp = []handlers.Error{
	{
		Code:    handlers.Forbidden,
		Message: "The provided authentication is not authorized to access this resource",
	},
}

func AuthMiddleware(
	identityDecoder interface {
		Decode(string) (*v1.Identity, error)
	},
) func(handler nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(handler nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
			if _, authenticated := ctx.Value(handlers.CookieAuthScopes).([]string); !authenticated {
				return handler(ctx, w, r, request)
			}

			cookie, err := r.Cookie("SESSIONID")
			if err != nil {
				slog.ErrorContext(ctx, "could not retrieve cookie", slog.String("cookieName", "SESSIONID"))
				WriteErrorResponse(ctx, w, http.StatusUnauthorized, newErrorResponse(ctx, UnauthorizedResp))
				return nil, nil
			}

			identity, err := identityDecoder.Decode(cookie.Value)
			if err != nil {
				if errors.Is(err, v1.ErrAuthorizationExpired) {
					WriteErrorResponse(ctx, w, http.StatusUnauthorized, newErrorResponse(ctx, AuthorizationExpiredResp))
					return nil, nil
				}
				slog.ErrorContext(ctx, "could not decode identity", slog.Any("error", err))
				WriteErrorResponse(ctx, w, http.StatusForbidden, newErrorResponse(ctx, ForbiddenResp))
				return nil, nil
			}

			ctx = context.WithValue(ctx, userCtxKey, identity)
			return handler(ctx, w, r, request)
		}
	}
}
