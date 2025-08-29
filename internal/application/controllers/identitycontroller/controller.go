package identitycontroller

import (
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/otel"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	identityEncoder IdentityEncoder
	db              application.Database
	tracer          trace.Tracer
}

func New(dep Dependencies) *Controller {
	traceProvider := otel.Provider.TracerProvider()
	return &Controller{
		identityEncoder: dep.IdentityEncoder,
		db:              dep.Database,
		tracer:          traceProvider.Tracer("identityController"),
	}
}
