package identitycontroller

import (
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/otel"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	identityEncoder IdentityEncoder
	db              ports.LocalDatabase
	tracer          trace.Tracer
}

func New(dep Dependencies) *Controller {
	traceProvider := otel.Provider.TracerProvider()

	return &Controller{
		identityEncoder: dep.IdentityEncoder,
		db:              dep.LocalDatabase,
		tracer:          traceProvider.Tracer("identityController"),
	}
}
