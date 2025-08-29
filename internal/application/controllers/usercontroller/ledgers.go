package usercontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/kset"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	LedgerController struct {
		db     application.Database
		tracer trace.Tracer
	}

	CreateLedgerRequest struct {
		Actor domain.ID
		Name  string
	}

	CreateLedgerResponse struct {
		ID domain.ID
	}

	CreateExpenseRequest struct {
		Actor          domain.ID
		LedgerID       domain.ID
		Name           string
		ExpenseDate    time.Time
		PendingRecords []domain.PendingRecord
	}

	CreateExpenseResponse struct {
		ID domain.ID
	}

	GetExpensesRequest struct {
		Actor    domain.ID
		LedgerID domain.ID
		Cursor   time.Time
		Limit    int32
	}

	GetExpensesResponse struct {
		Expenses []v1.LedgerExpenseSummary
		Cursor   *time.Time
	}

	AddMembersRequest struct {
		Actor    domain.ID
		LedgerID domain.ID
		Emails   []string
	}
)

func (c *LedgerController) Create(ctx context.Context, req CreateLedgerRequest) (resp *CreateLedgerResponse, err error) {
	ctx, span := c.tracer.Start(ctx, "create",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
		),
	)
	defer span.End()

	slog.Debug(ctx, "creating ledger", slog.With("req", req))

	err = c.db.Transaction(ctx, func(db application.Repositories) error {
		user, err := db.User().Find(ctx, req.Actor)
		if err != nil {
			return err
		}

		ledger, err := user.CreateLedger(req.Name)
		if err != nil {
			return err
		}

		if err := db.Ledger().Create(ctx, ledger); err != nil {
			return err
		}

		resp = &CreateLedgerResponse{
			ID: ledger.ID,
		}

		return db.User().Save(ctx, user)
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "creating ledger", err)
	}

	slog.Info(ctx, "ledger created", slog.WithStringer("ledger_id", resp.ID))

	return
}

func (c *LedgerController) CreateExpense(ctx context.Context, req CreateExpenseRequest) (resp *CreateExpenseResponse, err error) {
	ctx, span := c.tracer.Start(ctx, "createExpense",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	err = c.db.Transaction(ctx, func(db application.Repositories) error {
		ledger, err := db.Ledger().Find(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("finding ledger: %w", err)
		}

		expense, err := ledger.CreateExpense(domain.CreateExpenseRequest{
			Actor:          req.Actor,
			Name:           req.Name,
			ExpenseDate:    req.ExpenseDate,
			PendingRecords: req.PendingRecords,
		})
		if err != nil {
			return fmt.Errorf("creating expense: %w", err)
		}

		if err = db.Expense().Create(ctx, req.LedgerID, expense); err != nil {
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
		return nil, slog.ErrorReturn(ctx, "creating expense", err)
	}
	slog.Info(ctx, "expense created")
	return
}

type FindExpenseRequest struct {
	Actor     domain.ID
	LedgerID  domain.ID
	ExpenseID domain.ID
}

func (c *LedgerController) FindExpense(ctx context.Context, req FindExpenseRequest) (*domain.Expense, error) {
	ctx, span := c.tracer.Start(ctx, "findExpense",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
			attribute.Stringer("expense_id", req.ExpenseID),
		),
	)
	defer span.End()

	expense, err := c.db.Expense().Find(ctx, req.ExpenseID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger expense", err)
	}

	slog.Info(ctx, "ledger expense retrieved")

	return expense, nil
}

func (c *LedgerController) GetExpenses(ctx context.Context, req GetExpensesRequest) (*GetExpensesResponse, error) {
	ctx, span := c.tracer.Start(ctx, "getExpenses",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	req.Limit = max(1, req.Limit)

	expenses, err := c.db.Expense().GetByLedger(ctx, req.LedgerID, req.Cursor, req.Limit+1)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger expenses", err)
	}

	slog.Info(ctx, "ledger expenses retrieved")

	if len(expenses) == 0 {
		return &GetExpensesResponse{}, nil
	}

	var cursor *time.Time
	if len(expenses) == int(req.Limit)+1 {
		expenses = expenses[:len(expenses)-1]
		cursor = &expenses[len(expenses)-1].CreatedAt
	}

	return &GetExpensesResponse{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}

func (c *LedgerController) GetByIdentity(ctx context.Context, actor domain.ID) ([]domain.Ledger, error) {
	ctx, span := c.tracer.Start(ctx, "getByIdentity",
		trace.WithAttributes(
			attribute.Stringer("actor_id", actor),
		),
	)
	defer span.End()

	ledgers, err := c.db.Ledger().GetByUser(ctx, actor)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to list ledgers", err)
	}

	return ledgers, nil
}

type FindLedgerRequest struct {
	Actor    domain.ID
	LedgerID domain.ID
}

func (c *LedgerController) Find(ctx context.Context, req FindLedgerRequest) (*domain.Ledger, error) {
	ctx, span := c.tracer.Start(ctx, "find",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	ledger, err := c.db.Ledger().Find(ctx, req.LedgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to list ledgers", err)
	}

	return ledger, nil
}

// TODO(invitations): Here it's a simplification of the user membership process.
// We can always invert the flow and create invitation links, so the users click themselves
// We can also send invites through the system and they accept the invite through the API.
func (c *LedgerController) AddMembers(ctx context.Context, req AddMembersRequest) error {
	ctx, span := c.tracer.Start(ctx, "addMembers",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	transaction := func(db application.Repositories) error {
		users, err := db.User().ListByEmail(ctx, req.Emails)
		if err != nil {
			return slog.ErrorReturn(ctx, "failed to get users", err)
		}

		ledger, err := db.Ledger().Find(ctx, req.LedgerID)
		if err != nil {
			return slog.ErrorReturn(ctx, "failed to find ledger", err)
		}

		err = ledger.AddMember(req.Actor, kset.Select(func(u domain.User) domain.ID { return u.ID }, users...)...)
		if err != nil {
			return slog.ErrorReturn(ctx, "adding members", err)
		}

		if err := db.Ledger().Update(ctx, ledger); err != nil {
			return err
		}

		return nil
	}

	switch err := c.db.Transaction(ctx, transaction); {
	case err == nil:
		slog.Info(ctx, "added users to ledger")
		return nil
	default:
		return slog.ErrorReturn(ctx, "failed to add users to ledger", err)
	}
}

type GetMembersRequest struct {
	Actor    domain.ID
	LedgerID domain.ID
}

func (c *LedgerController) GetMembers(ctx context.Context, req GetMembersRequest) (map[domain.ID]*domain.LedgerMember, error) {
	ctx, span := c.tracer.Start(ctx, "getMembers",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
		),
	)
	defer span.End()

	members, err := c.db.Ledger().Find(ctx, req.LedgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger members balances", err)
	}

	slog.Info(ctx, "ledger members balances retrieved")

	return members.Members, nil
}

type CreateExpenseRecordRequest struct {
	Actor          domain.ID
	LedgerID       domain.ID
	ExpenseID      domain.ID
	PendingRecords []domain.PendingRecord
}

func (c *LedgerController) CreateExpenseRecord(ctx context.Context, req CreateExpenseRecordRequest) (resp *domain.Expense, err error) {
	ctx, span := c.tracer.Start(ctx, "createExpenseRecord",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.Actor),
			attribute.Stringer("ledger_id", req.LedgerID),
			attribute.Stringer("expense_id", req.ExpenseID),
		),
	)
	defer span.End()

	err = c.db.Transaction(ctx, func(db application.Repositories) error {
		ledger, err := db.Ledger().Find(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("fetching ledger: %w", err)
		}

		expense, err := db.Expense().Find(ctx, req.ExpenseID)
		if err != nil {
			return fmt.Errorf("fetching expense: %w", err)
		}

		if err := expense.CreateRecords(req.Actor, ledger, req.PendingRecords...); err != nil {
			return fmt.Errorf("appending new records: %w", err)
		}

		if err := db.Ledger().Update(ctx, ledger); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		if err := db.Expense().Update(ctx, expense); err != nil {
			return fmt.Errorf("updating expense: %w", err)
		}

		resp = expense
		return nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "creating expense record", err)
	}

	slog.Info(ctx, "expense records created")

	return
}

type DeleteExpenseRecordRequest struct {
	ActorID   domain.ID
	LedgerID  domain.ID
	ExpenseID domain.ID
	RecordID  domain.ID
}

func (c *LedgerController) DeleteExpenseRecord(ctx context.Context, req DeleteExpenseRecordRequest) error {
	ctx, span := c.tracer.Start(ctx, "createExpenseRecord",
		trace.WithAttributes(
			attribute.Stringer("actor_id", req.ActorID),
			attribute.Stringer("ledger_id", req.LedgerID),
			attribute.Stringer("expense_id", req.ExpenseID),
			attribute.Stringer("record_id", req.RecordID),
		),
	)
	defer span.End()

	err := c.db.Transaction(ctx, func(db application.Repositories) error {
		expense, err := db.Expense().Find(ctx, req.ExpenseID)
		if err != nil {
			return fmt.Errorf("finding expense: %w", err)
		}

		ledger, err := db.Ledger().Find(ctx, req.LedgerID)
		if err != nil {
			return fmt.Errorf("finding ledger: %w", err)
		}

		if err = expense.DeleteRecord(req.ActorID, ledger, req.RecordID); err != nil {
			return fmt.Errorf("deleting record: %w", err)
		}

		if err = db.Ledger().Update(ctx, ledger); err != nil {
			return fmt.Errorf("updating ledger: %w", err)
		}

		if err = db.Expense().Update(ctx, expense); err != nil {
			return fmt.Errorf("updating expense: %w", err)
		}

		return nil
	})
	if err != nil {
		return slog.ErrorReturn(ctx, "deleting expense record", err)
	}

	slog.Info(ctx, "ledger expense record deleted")

	return nil
}
