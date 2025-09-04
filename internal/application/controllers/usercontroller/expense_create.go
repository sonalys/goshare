package usercontroller

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/ports"
	"github.com/sonalys/goshare/pkg/slog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	CreateExpenseRequest struct {
		ActorID        domain.ID
		LedgerID       domain.ID
		Name           string
		ExpenseDate    time.Time
		PendingRecords []domain.PendingRecord
	}

	CreateExpenseResponse struct {
		ID domain.ID
	}
)

func (c *expenseController) Create(ctx context.Context, req CreateExpenseRequest) (resp *CreateExpenseResponse, err error) {
	ctx, span := c.tracer.Start(ctx, "create",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	slog.Debug(ctx, "creating expense", slog.With("req", req))

	err = c.db.Transaction(ctx, func(db ports.LocalRepositories) error {
		ledger, err := db.Ledger().Get(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("finding ledger: %w", err)
		}

		if !ledger.CanManageExpenses(req.ActorID) {
			return fmt.Errorf("authorizing user ledger expense management: %w", v1.ErrForbidden)
		}

		expense, err := ledger.CreateExpense(domain.CreateExpenseRequest{
			Creator:        req.ActorID,
			Name:           req.Name,
			ExpenseDate:    req.ExpenseDate,
			PendingRecords: req.PendingRecords,
		})
		if err != nil {
			return fmt.Errorf("creating expense: %w", err)
		}

		if err = db.Expense().Create(ctx, expense); err != nil {
			return fmt.Errorf("saving expense: %w", err)
		}

		if err = db.Ledger().Update(ctx, ledger); err != nil {
			return fmt.Errorf("saving ledger: %w", err)
		}

		resp = &CreateExpenseResponse{
			ID: expense.ID,
		}

		return nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "committing transaction", err)
	}

	slog.Info(ctx, "expense created")

	return
}
