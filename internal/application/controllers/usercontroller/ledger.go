package usercontroller

import (
	"github.com/sonalys/goshare/internal/application"
	"go.opentelemetry.io/otel/trace"
)

type (
	ledgerController struct {
		db     application.Database
		tracer trace.Tracer
	}
)
