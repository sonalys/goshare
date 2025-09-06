package middlewares

import (
	"context"
	"errors"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

type (
	contextKey string

	IdentityDecoder interface {
		Decode(jwt string) (*application.Identity, error)
	}

	SecurityHandler struct {
		controller IdentityDecoder
	}
)

const identityContextKey = contextKey("identity-key")

func NewSecurityHandler(c IdentityDecoder) *SecurityHandler {
	return &SecurityHandler{
		controller: c,
	}
}

func (h *SecurityHandler) GetIdentity(ctx context.Context) (*application.Identity, error) {
	identity, ok := ctx.Value(identityContextKey).(*application.Identity)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	return identity, nil
}

func (h *SecurityHandler) HandleCookieAuth(ctx context.Context, operationName server.OperationName, t server.CookieAuth) (context.Context, error) {
	identity, err := h.controller.Decode(t.APIKey)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, identityContextKey, identity), nil
}
