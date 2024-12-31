package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/sonalys/goshare/cmd/server/handlers"
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

const userCtxKey = "user-key"

type Identity struct {
	Email  string
	UserID uuid.UUID
}

func GetIdentity(ctx context.Context) (*Identity, error) {
	identity, ok := ctx.Value(userCtxKey).(*Identity)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return identity, nil
}

func AuthMiddleware(handler nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		if _, authenticated := ctx.Value(handlers.CookieAuthScopes).([]string); !authenticated {
			return handler(ctx, w, r, request)
		}

		cookie, err := r.Cookie("SESSIONID")
		if err != nil {
			slog.ErrorContext(ctx, "could not retrieve cookie", slog.String("cookieName", "SESSIONID"))
			return nil, err
		}

		var claims jwt.MapClaims
		token, err := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("my-secret-key"), nil
		})

		if err != nil || !token.Valid {
			slog.ErrorContext(ctx, "token is not valid", slog.Any("error", err), slog.Bool("isValid", token.Valid))
			return nil, err
		}

		identity := &Identity{
			Email:  claims["email"].(string),
			UserID: uuid.MustParse(claims["userID"].(string)),
		}

		ctx = context.WithValue(ctx, userCtxKey, identity)
		return handler(ctx, w, r, request)
	}
}
