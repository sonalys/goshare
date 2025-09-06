package identitycontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/otel"
	"go.opentelemetry.io/otel/trace"
)

type (
	Controller interface {
		Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
		Register(ctx context.Context, req RegisterRequest) (resp *RegisterResponse, err error)
	}

	controller struct {
		identityEncoder IdentityEncoder
		db              ports.LocalDatabase
		tracer          trace.Tracer
	}
)

func New(dep Dependencies) Controller {
	traceProvider := otel.Provider.TracerProvider()

	return &controller{
		identityEncoder: dep.IdentityEncoder,
		db:              dep.LocalDatabase,
		tracer:          traceProvider.Tracer("identityController"),
	}
}
