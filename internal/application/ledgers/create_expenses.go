package ledgers

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
)

type (
	CreateExpenseRequest struct {
		UserID       uuid.UUID
		LedgerID     uuid.UUID
		CategoryID   *uuid.UUID
		Amount       int32
		Name         string
		ExpenseDate  time.Time
		UserBalances []v1.ExpenseUserBalance
	}

	CreateExpenseResponse struct {
		ID uuid.UUID
	}
)

func (r CreateExpenseRequest) Validate() error {
	var errs v1.FormError

	if r.LedgerID == uuid.Nil {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("ledger_id"))
	}

	if r.Amount <= 0 {
		errs.Fields = append(errs.Fields, v1.NewFieldRangeError("amount", 0, math.MaxInt32))
	}

	if r.Name == "" {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("name"))
	}

	if r.ExpenseDate.IsZero() {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("expense_date"))
	}

	if len(r.UserBalances) == 0 {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("user_balances"))
	}

	var balanceSum int32
	var totalPaid int32
	for i, ub := range r.UserBalances {
		balanceSum += ub.Balance

		if ub.Balance > 0 {
			totalPaid += ub.Balance
		}

		if ub.UserID == uuid.Nil {
			errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("user_balances["+fmt.Sprint(i)+"].user_id"))
		}
	}

	if balanceSum != 0 {
		errs.Fields = append(errs.Fields, v1.FieldError{
			Field: "user_balances",
			Cause: fmt.Errorf("%w: sum should be equal to 0. got %s", v1.ErrInvalidValue, v1.NewMoney(balanceSum, 2, "$")),
		})
	}

	if totalPaid != r.Amount {
		errs.Fields = append(errs.Fields, v1.FieldError{
			Field: "user_balances",
			Cause: fmt.Errorf("%w: total paid balance should match expense amount. expected %s, got %s", v1.ErrInvalidValue, v1.NewMoney(r.Amount, 2, "$"), v1.NewMoney(totalPaid, 2, "$")),
		})
	}

	return errs.Validate()
}

func (c *Controller) CreateExpense(ctx context.Context, req CreateExpenseRequest) (*CreateExpenseResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.CreateExpense")
	defer span.End()

	if err := req.Validate(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "invalid request", slog.Any("error", err))
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	expense := &v1.Expense{
		ID:           uuid.New(),
		CategoryID:   req.CategoryID,
		LedgerID:     req.LedgerID,
		Amount:       req.Amount,
		Name:         req.Name,
		ExpenseDate:  req.ExpenseDate,
		UserBalances: req.UserBalances,
		CreatedAt:    time.Now(),
		CreatedBy:    req.UserID,
		UpdatedAt:    time.Now(),
		UpdatedBy:    req.UserID,
	}

	if err := c.expenseRepository.Create(ctx, expense); err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to create expense", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	slog.InfoContext(ctx, "expense created", slog.String("expense_id", expense.ID.String()))

	return &CreateExpenseResponse{
		ID: expense.ID,
	}, nil
}
