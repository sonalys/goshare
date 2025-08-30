package usercontroller

import (
	"github.com/sonalys/goshare/internal/ports"
	"go.opentelemetry.io/otel/trace"
)

type expenseController struct {
	db     ports.LocalDatabase
	tracer trace.Tracer
}
