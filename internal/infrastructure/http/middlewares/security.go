package middlewares

import (
	"context"
	"errors"

	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	v1 "github.com/sonalys/goshare/pkg/v1"
)

type (
	contextKey string

	IdentityDecoder interface {
		Decode(jwt string) (*v1.Identity, error)
	}

	SecurityHandler struct {
		controller IdentityDecoder
	}
)

const identityContextKey = contextKey("identity-key")

func GetIdentity(ctx context.Context) (*v1.Identity, error) {
	identity, ok := ctx.Value(identityContextKey).(*v1.Identity)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	return identity, nil
}

func NewSecurityHandler(c IdentityDecoder) server.SecurityHandler {
	return &SecurityHandler{
		controller: c,
	}
}

func (h *SecurityHandler) HandleCookieAuth(ctx context.Context, operationName server.OperationName, t server.CookieAuth) (context.Context, error) {
	identity, err := h.controller.Decode(t.APIKey)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, identityContextKey, identity), nil
}
