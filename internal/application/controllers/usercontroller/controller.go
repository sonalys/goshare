package usercontroller

import (
	"github.com/sonalys/goshare/internal/application/pkg/otel"
)

type (
	Controller struct {
		*ledgerController
		*recordsController
		*expenseController
	}
)

func New(dep Dependencies) *Controller {
	traceProvider := otel.Provider.TracerProvider()
	return &Controller{
		ledgerController: &ledgerController{
			db:     dep.Database,
			tracer: traceProvider.Tracer("userController.ledger"),
		},
		recordsController: &recordsController{
			db:     dep.Database,
			tracer: traceProvider.Tracer("userController.record"),
		},
		expenseController: &expenseController{
			db:     dep.Database,
			tracer: traceProvider.Tracer("userController.expense"),
		},
	}
}

func (c *Controller) Ledgers() *ledgerController {
	return c.ledgerController
}

func (c *Controller) Records() *recordsController {
	return c.recordsController
}

func (c *Controller) Expenses() *expenseController {
	return c.expenseController
}
