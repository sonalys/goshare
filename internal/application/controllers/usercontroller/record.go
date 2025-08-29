package usercontroller

import (
	"github.com/sonalys/goshare/internal/application"
	"go.opentelemetry.io/otel/trace"
)

type (
	recordsController struct {
		db     application.Database
		tracer trace.Tracer
	}
)
