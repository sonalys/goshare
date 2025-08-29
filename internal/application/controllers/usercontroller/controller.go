package usercontroller

import (
	"github.com/sonalys/goshare/internal/application/pkg/otel"
)

type (
	Controller struct {
		*LedgerController
	}
)

func New(dep Dependencies) *Controller {
	traceProvider := otel.Provider.TracerProvider()
	return &Controller{
		LedgerController: &LedgerController{
			db:     dep.Database,
			tracer: traceProvider.Tracer("userController.ledgers"),
		},
	}
}

func (c *Controller) Ledgers() *LedgerController {
	return c.LedgerController
}
