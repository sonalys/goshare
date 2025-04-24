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

	GetExpensesParams struct {
		UserID   v1.ID
		LedgerID v1.ID
		Cursor   time.Time
		Limit    int32
	}

	GetExpensesResult struct {
		Expenses []v1.LedgerExpenseSummary
		Cursor   *time.Time
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
		return nil, err
	}

	var totalAmount int32

	for _, record := range req.Records {
		if record.Type == v1.RecordTypeDebt {
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

	switch err := c.expenseRepository.Create(ctx, req.LedgerID, func(ledger *v1.Ledger) (*v1.Expense, error) {
		if !ledger.IsParticipant(req.UserID) {
			return nil, v1.ErrUserNotAMember
		}
		return expense, nil
	}); {
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
		return nil, err
	default:
		slog.InfoContext(ctx, "expense created", slog.String("expense_id", expense.ID.String()))

		return &CreateExpenseResponse{
			ID: expense.ID,
		}, nil
	}
}

func (c *Controller) FindExpense(ctx context.Context, _ v1.ID, expenseID v1.ID) (*v1.Expense, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.FindExpense")
	defer span.End()

	logFields := []any{
		slog.String("expense_id", expenseID.String()),
	}

	expense, err := c.expenseRepository.Find(ctx, expenseID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get ledger expense", append(logFields, slog.Any("error", err))...)
		return nil, err
	}

	slog.InfoContext(ctx, "ledger expense retrieved")

	return expense, nil
}

func (c *Controller) GetExpenses(ctx context.Context, params GetExpensesParams) (*GetExpensesResult, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetExpenses")
	defer span.End()

	params.Limit = max(1, params.Limit)

	logFields := []any{
		slog.String("ledger_id", params.LedgerID.String()),
	}

	ledger, err := c.ledgerRepository.Find(ctx, params.LedgerID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get ledger", append(logFields, slog.Any("error", err))...)
		return nil, err
	}

	if !ledger.IsParticipant(params.UserID) {
		return nil, v1.ErrUserNotAMember
	}

	expenses, err := c.expenseRepository.GetByLedger(ctx, params.LedgerID, params.Cursor, params.Limit+1)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get ledger expenses", append(logFields, slog.Any("error", err))...)
		return nil, err
	}

	if len(expenses) == 0 {
		return nil, nil
	}

	slog.InfoContext(ctx, "ledger expenses retrieved", logFields...)

	var cursor *time.Time
	if len(expenses) == int(params.Limit)+1 {
		expenses = expenses[:len(expenses)-1]
		cursor = &expenses[len(expenses)-1].CreatedAt
	}

	return &GetExpensesResult{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}
