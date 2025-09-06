package ledgers_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers"
	"github.com/sonalys/goshare/internal/mocks/application/controllers/usercontrollermock"
	"github.com/sonalys/goshare/internal/mocks/portsmock"
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
		SecurityHandler: &portsmock.SecurityHandler{},
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
