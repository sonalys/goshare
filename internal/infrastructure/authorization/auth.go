package authorization

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/application/usecases"
)

type Authorizer struct {
	controllers.LedgerRepository
	controllers.UserRepository
}

func (a *Authorizer) Authorize(ctx context.Context, action usecases.AuthorizationActionType, resource usecases.AuthorizationResource, subject usecases.AuthorizationResource) error {
	switch action {
	case usecases.ActionLedgerExpenseWrite:
		return nil
	default:
		return v1.ErrForbidden
	}
}

var _ usecases.Authorizer = &Authorizer{}
