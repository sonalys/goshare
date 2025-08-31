package ledgers_test

import (
	"context"
	"testing"
	"time"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers"
	"github.com/sonalys/goshare/internal/mocks/application/controllers/usercontrollermock"
	"github.com/sonalys/goshare/internal/mocks/portsmock"
	v1 "github.com/sonalys/goshare/pkg/v1"
	"github.com/stretchr/testify/assert"
)

type testSetup struct {
	*portsmock.SecurityHandler
	*usercontrollermock.Controller
	*usercontrollermock.LedgerController
	*usercontrollermock.ExpenseController
	*usercontrollermock.RecordsController
}

func withIdentity(identity *v1.Identity) func(*testSetup) {
	return func(ts *testSetup) {
		ts.GetIdentityFunc = func(ctx context.Context) (*v1.Identity, error) {
			return identity, nil
		}
	}
}

func unauthenticated() func(*testSetup) {
	return func(ts *testSetup) {
		ts.GetIdentityFunc = func(ctx context.Context) (*v1.Identity, error) {
			return nil, assert.AnError
		}
	}
}

func setup(_ *testing.T, hooks ...func(*testSetup)) (*ledgers.Router, *testSetup) {
	ledgerController := &usercontrollermock.LedgerController{}
	expenseController := &usercontrollermock.ExpenseController{}
	recordsController := &usercontrollermock.RecordsController{}

	ts := &testSetup{
		SecurityHandler: &portsmock.SecurityHandler{
			GetIdentityFunc: func(ctx context.Context) (*v1.Identity, error) {
				return &v1.Identity{
					Email:  "email@example.com",
					UserID: domain.NewID(),
					Exp:    domain.Now().Add(time.Hour).Unix(),
				}, nil
			},
		},
		Controller: &usercontrollermock.Controller{
			ExpensesFunc: func() usercontroller.ExpenseController { return expenseController },
			LedgersFunc:  func() usercontroller.LedgerController { return ledgerController },
			RecordsFunc:  func() usercontroller.RecordsController { return recordsController },
		},
		LedgerController:  ledgerController,
		ExpenseController: expenseController,
		RecordsController: recordsController,
	}

	for _, hook := range hooks {
		hook(ts)
	}

	return ledgers.New(ts.SecurityHandler, ts.Controller), ts
}
