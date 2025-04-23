package ledgers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	CreateExpenseRequest struct {
		UserID      v1.ID
		LedgerID    v1.ID
		Name        string
		ExpenseDate time.Time
		Records     []v1.Record
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

	if r.Name == "" {
		errs = append(errs, v1.NewRequiredFieldError("name"))
	}

	if r.ExpenseDate.IsZero() {
		errs = append(errs, v1.NewRequiredFieldError("expense_date"))
	}

	if len(r.Records) == 0 {
		errs = append(errs, v1.NewRequiredFieldError("user_balances"))
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

	var totalAmount int32

	for _, record := range req.Records {
		if record.Type == v1.RecordTypeExpense {
			totalAmount += record.Amount
		}
	}

	expense := &v1.Expense{
		ID:          v1.NewID(),
		LedgerID:    req.LedgerID,
		Name:        req.Name,
		Amount:      totalAmount,
		ExpenseDate: req.ExpenseDate,
		Records:     req.Records,
		CreatedAt:   time.Now(),
		CreatedBy:   req.UserID,
		UpdatedAt:   time.Now(),
		UpdatedBy:   req.UserID,
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
