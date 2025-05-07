package authorization

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/application/usecases"
	"github.com/sonalys/goshare/internal/domain"
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
		return domain.ErrForbidden
	}
}

var _ usecases.Authorizer = &Authorizer{}
