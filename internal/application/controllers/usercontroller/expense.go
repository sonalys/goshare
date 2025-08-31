package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/ports"
	"go.opentelemetry.io/otel/trace"
)

type (
	ExpenseController interface {
		Create(ctx context.Context, req CreateExpenseRequest) (resp *CreateExpenseResponse, err error)
		Get(ctx context.Context, req GetExpenseRequest) (*domain.Expense, error)
		List(ctx context.Context, req ListExpensesRequest) (*ListExpensesResponse, error)
	}

	expenseController struct {
		db     ports.LocalDatabase
		tracer trace.Tracer
	}
)
