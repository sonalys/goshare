package testutils

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/infrastructure/http/router"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers/expenses"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/ledgers/expenses/records"
	"github.com/sonalys/goshare/internal/infrastructure/http/router/users"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
	"github.com/sonalys/goshare/internal/mocks/application/controllers/identitycontrollermock"
	"github.com/sonalys/goshare/internal/mocks/application/controllers/usercontrollermock"
	"github.com/sonalys/goshare/internal/mocks/portsmock"
	"github.com/stretchr/testify/assert"
)

type testSetup struct {
	SecurityHandler    *portsmock.SecurityHandler
	IdentityController *identitycontrollermock.Controller
	UserController     *usercontrollermock.Controller
	LedgerController   *usercontrollermock.LedgerController
	ExpenseController  *usercontrollermock.ExpenseController
	RecordsController  *usercontrollermock.RecordsController
}

func WithIdentity(identity *application.Identity) func(*testSetup) {
	return func(ts *testSetup) {
		ts.SecurityHandler.GetIdentityFunc = func(ctx context.Context) (*application.Identity, error) {
			if identity == nil {
				return nil, assert.AnError
			}

			return identity, nil
		}
	}
}

func Setup(_ *testing.T, hooks ...func(*testSetup)) (server.Handler, *testSetup) {
	ledgerController := &usercontrollermock.LedgerController{}
	expenseController := &usercontrollermock.ExpenseController{}
	recordsController := &usercontrollermock.RecordsController{}

	ts := &testSetup{
		SecurityHandler:    &portsmock.SecurityHandler{},
		IdentityController: &identitycontrollermock.Controller{},
		UserController: &usercontrollermock.Controller{
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

	return &router.Router{
		LedgersHandler:  ledgers.New(ts.SecurityHandler, ts.UserController),
		UsersHandler:    users.New(ts.SecurityHandler, ts.IdentityController, ts.UserController),
		ExpensesHandler: expenses.New(ts.SecurityHandler, ts.UserController),
		RecordsHandler:  records.New(ts.SecurityHandler, ts.UserController),
	}, ts
}
