package middlewares

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	v1 "github.com/sonalys/goshare/pkg/v1"
)

type contextKey string

type identityController interface {
	Decode(string) (*v1.Identity, error)
}

var identityContextKey contextKey = "identity-key"

func GetIdentity(ctx context.Context) (*v1.Identity, error) {
	identity, ok := ctx.Value(identityContextKey).(*v1.Identity)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return identity, nil
}

type SecurityHandler struct {
	controller identityController
}

func (h *SecurityHandler) HandleCookieAuth(ctx context.Context, operationName server.OperationName, t server.CookieAuth) (context.Context, error) {
	identity, err := h.controller.Decode(t.APIKey)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, identityContextKey, identity), nil
}

func NewSecurityHandler(c identityController) *SecurityHandler {
	return &SecurityHandler{
		controller: c,
	}
}
