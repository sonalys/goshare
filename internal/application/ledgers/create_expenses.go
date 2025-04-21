package ledgers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	CreateExpenseRequest struct {
		UserID       v1.ID
		LedgerID     v1.ID
		CategoryID   *v1.ID
		Amount       int32
		Name         string
		ExpenseDate  time.Time
		UserBalances []v1.ExpenseUserBalance
	}

	CreateExpenseResponse struct {
		ID v1.ID
	}
)

func (r CreateExpenseRequest) Validate() error {
	var errs v1.FormError

	if r.LedgerID.IsEmpty() {
		errs = append(errs, v1.NewRequiredFieldError("ledger_id"))
	}

	if r.Amount <= 0 {
		errs = append(errs, v1.NewFieldRangeError("amount", 0, math.MaxInt32))
	}

	if r.Name == "" {
		errs = append(errs, v1.NewRequiredFieldError("name"))
	}

	if r.ExpenseDate.IsZero() {
		errs = append(errs, v1.NewRequiredFieldError("expense_date"))
	}

	if len(r.UserBalances) == 0 {
		errs = append(errs, v1.NewRequiredFieldError("user_balances"))
	}

	var balanceSum int32
	var totalPaid int32
	for i, ub := range r.UserBalances {
		balanceSum += ub.Balance

		if ub.Balance > 0 {
			totalPaid += ub.Balance
		}

		if ub.UserID.IsEmpty() {
			errs = append(errs, v1.NewRequiredFieldError("user_balances["+fmt.Sprint(i)+"].user_id"))
		}
	}

	if balanceSum != 0 {
		errs = append(errs, v1.FieldError{
			Field: "user_balances",
			Cause: fmt.Errorf("%w: sum should be equal to 0. got %s", v1.ErrInvalidValue, v1.NewMoney(balanceSum, 2, "$")),
		})
	}

	if totalPaid > r.Amount {
		errs = append(errs, v1.FieldError{
			Field: "user_balances",
			Cause: fmt.Errorf("%w: total payment should be less or equal to the expense value. got %s", v1.ErrInvalidValue, v1.NewMoney(totalPaid, 2, "$")),
		})
	}

	return errs.Validate()
}

func (c *Controller) CreateExpense(ctx context.Context, req CreateExpenseRequest) (*CreateExpenseResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.CreateExpense")
	defer span.End()

	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, "invalid request", slog.Any("error", err))
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	expense := &v1.Expense{
		ID:           v1.NewID(),
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

	switch err := c.expenseRepository.Create(ctx, expense); {
	case errors.Is(err, v1.ErrUserNotAMember):
		if fieldErr := new(v1.FieldError); errors.As(err, fieldErr) {
			return nil, v1.FieldError{
				Cause:    v1.ErrUserNotAMember,
				Field:    fmt.Sprintf("user_balances.%d.user_id", fieldErr.Metadata.Index),
				Metadata: fieldErr.Metadata,
			}
		}
		return nil, err
	case err != nil:
		slog.ErrorContext(ctx, "failed to create expense", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create expense: %w", err)
	default:
		slog.InfoContext(ctx, "expense created", slog.String("expense_id", expense.ID.String()))

		return &CreateExpenseResponse{
			ID: expense.ID,
		}, nil
	}
}
