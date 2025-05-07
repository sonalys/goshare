package usecases

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
)

type (
	AuthorizationResource interface {
		ID() domain.ID
		Type() string
	}

	resource struct {
		t  string
		id domain.ID
	}

	AuthorizationActionType int

	Authorizer interface {
		Authorize(ctx context.Context, action AuthorizationActionType, resource AuthorizationResource, subject AuthorizationResource) error
	}
)

const (
	resourceTypeUnspecified AuthorizationActionType = iota
	ResourceTypeLedger
	ResourceTypeUser

	actionUnspecified AuthorizationActionType = iota
	ActionLedgerExpenseCreate
	ActionLedgerExpenseRead
	ActionLedgerExpenseWrite
)

func ResourceLedger(id domain.ID) AuthorizationResource {
	return resource{
		t:  "ledger",
		id: id,
	}
}

func ResourceUser(id domain.ID) AuthorizationResource {
	return resource{
		t:  "user",
		id: id,
	}
}

func (r resource) ID() domain.ID {
	return r.id
}

func (r resource) Type() string {
	return r.t
}
