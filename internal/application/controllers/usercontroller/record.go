package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/ports"
	"go.opentelemetry.io/otel/trace"
)

type (
	RecordsController interface {
		Create(ctx context.Context, req CreateExpenseRecordRequest) (resp *domain.Expense, err error)
		Delete(ctx context.Context, req DeleteExpenseRecordRequest) error
	}

	recordsController struct {
		db     ports.LocalDatabase
		tracer trace.Tracer
	}
)
