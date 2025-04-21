package api

import (
	"context"
	"fmt"

	"github.com/sonalys/goshare/cmd/server/handlers"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

const userCtxKey = "user-key"

func getIdentity(ctx context.Context) (*v1.Identity, error) {
	identity, ok := ctx.Value(userCtxKey).(*v1.Identity)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return identity, nil
}

type SecurityHandler struct {
	identityDecoder interface {
		Decode(string) (*v1.Identity, error)
	}
}

func (h *SecurityHandler) HandleCookieAuth(ctx context.Context, operationName handlers.OperationName, t handlers.CookieAuth) (context.Context, error) {
	identity, err := h.identityDecoder.Decode(t.APIKey)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, userCtxKey, identity)
	return ctx, nil
}

func NewSecurityHandler(identityDecoder interface {
	Decode(string) (*v1.Identity, error)
}) *SecurityHandler {
	return &SecurityHandler{
		identityDecoder: identityDecoder,
	}
}
