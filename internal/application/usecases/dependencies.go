package usecases

import (
	"context"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

type (
	AuthorizationResource interface {
		ID() v1.ID
		Type() string
	}

	resource struct {
		t  string
		id v1.ID
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

func ResourceLedger(id v1.ID) AuthorizationResource {
	return resource{
		t:  "ledger",
		id: id,
	}
}

func ResourceUser(id v1.ID) AuthorizationResource {
	return resource{
		t:  "user",
		id: id,
	}
}

func (r resource) ID() v1.ID {
	return r.id
}

func (r resource) Type() string {
	return r.t
}
