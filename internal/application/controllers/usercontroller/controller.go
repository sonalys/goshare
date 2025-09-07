package usercontroller

import (
	"github.com/sonalys/goshare/internal/pkg/otel"
)

type (
	Controller interface {
		Expenses() ExpenseController
		Ledgers() LedgerController
		Records() RecordsController
	}

	controller struct {
		*ledgerController
		*recordsController
		*expenseController
	}
)

func New(dep Dependencies) Controller {
	traceProvider := otel.Provider.TracerProvider()

	return &controller{
		ledgerController: &ledgerController{
			db:     dep.LocalDatabase,
			tracer: traceProvider.Tracer("userController.ledger"),
		},
		recordsController: &recordsController{
			db:     dep.LocalDatabase,
			tracer: traceProvider.Tracer("userController.record"),
		},
		expenseController: &expenseController{
			db:     dep.LocalDatabase,
			tracer: traceProvider.Tracer("userController.expense"),
		},
	}
}

func (c *controller) Ledgers() LedgerController {
	return c.ledgerController
}

func (c *controller) Records() RecordsController {
	return c.recordsController
}

func (c *controller) Expenses() ExpenseController {
	return c.expenseController
}
