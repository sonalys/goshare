package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/goshare/internal/application/pkg/otel"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/kset"
)

type (
	Ledgers struct {
		db Database
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

func (c *Ledgers) Create(ctx context.Context, req CreateLedgerRequest) (resp *CreateLedgerResponse, err error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.Create")
	defer span.End()

	ctx = slog.Context(ctx, slog.WithStringer("actor", req.Actor))

	slog.Debug(ctx, "creating ledger", slog.WithAny("req", req))

	err = c.db.Transaction(ctx, func(db Database) error {
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

func (c *Ledgers) CreateExpense(ctx context.Context, req CreateExpenseRequest) (resp *CreateExpenseResponse, err error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.CreateExpense")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("actor", req.Actor),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	err = c.db.Transaction(ctx, func(db Database) error {
		ledger, err := db.Ledger().Find(ctx, req.LedgerID)
		if err != nil {
			return err
		}

		expense, err := ledger.CreateExpense(domain.CreateExpenseRequest{
			Actor:          req.Actor,
			Name:           req.Name,
			ExpenseDate:    req.ExpenseDate,
			PendingRecords: req.PendingRecords,
		})
		if err != nil {
			return err
		}

		err = db.Expense().Create(ctx, req.LedgerID, expense)
		if err != nil {
			return err
		}

		err = db.Ledger().Update(ctx, ledger)
		if err != nil {
			return err
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

func (c *Ledgers) FindExpense(ctx context.Context, ledgerID domain.ID, expenseID domain.ID) (*domain.Expense, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.FindExpense")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("ledger_id", ledgerID),
		slog.WithStringer("expense_id", expenseID),
	)

	expense, err := c.db.Expense().Find(ctx, expenseID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger expense", err)
	}

	slog.Info(ctx, "ledger expense retrieved")

	return expense, nil
}

func (c *Ledgers) GetExpenses(ctx context.Context, req GetExpensesRequest) (*GetExpensesResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetExpenses")
	defer span.End()

	req.Limit = max(1, req.Limit)

	ctx = slog.Context(ctx,
		slog.WithStringer("actor", req.Actor),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

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

func (c *Ledgers) GetByIdentity(ctx context.Context, identity domain.ID) ([]domain.Ledger, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.ListByUser")
	defer span.End()

	ctx = slog.Context(ctx, slog.WithStringer("actor", identity))

	ledgers, err := c.db.Ledger().GetByUser(ctx, identity)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to list ledgers", err)
	}

	return ledgers, nil
}

func (c *Ledgers) Find(ctx context.Context, ledgerID domain.ID) (*domain.Ledger, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.Find")
	defer span.End()

	ctx = slog.Context(ctx, slog.WithStringer("ledgerID", ledgerID))

	ledger, err := c.db.Ledger().Find(ctx, ledgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to list ledgers", err)
	}

	return ledger, nil
}

// TODO(invitations): Here it's a simplification of the user membership process.
// We can always invert the flow and create invitation links, so the users click themselves
// We can also send invites through the system and they accept the invite through the API.
func (c *Ledgers) AddMembers(ctx context.Context, req AddMembersRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.AddMembers")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("actor", req.Actor),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	transaction := func(db Database) error {
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

func (c *Ledgers) GetMembers(ctx context.Context, ledgerID domain.ID) (map[domain.ID]*domain.LedgerMember, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetMembers")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("ledger_id", ledgerID),
	)

	members, err := c.db.Ledger().Find(ctx, ledgerID)
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

func (c *Ledgers) CreateExpenseRecord(ctx context.Context, req CreateExpenseRecordRequest) (resp *domain.Expense, err error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.CreateExpenseRecord")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("actor", req.Actor),
		slog.WithStringer("expenseID", req.ExpenseID),
		slog.WithStringer("ledgerID", req.LedgerID),
	)

	err = c.db.Transaction(ctx, func(db Database) error {
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

func (c *Ledgers) DeleteExpenseRecord(ctx context.Context, req DeleteExpenseRecordRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.CreateExpenseRecord")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("actor", req.ActorID),
		slog.WithStringer("ledgerID", req.LedgerID),
		slog.WithStringer("expenseID", req.ExpenseID),
		slog.WithStringer("recordID", req.RecordID),
	)

	err := c.db.Transaction(ctx, func(db Database) error {
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
