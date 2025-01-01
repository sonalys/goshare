package ledgers

import (
	"context"
	"fmt"
	"log/slog"
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
		CategoryID   uuid.UUID
		Amount       int32
		Name         string
		ExpenseDate  time.Time
		UserBalances []v1.ExpenseUserBalance
	}

	CreateExpenseResponse struct {
		ID uuid.UUID
	}
)

func (c *Controller) CreateExpense(ctx context.Context, req CreateExpenseRequest) (*CreateExpenseResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.CreateExpense")
	defer span.End()

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
