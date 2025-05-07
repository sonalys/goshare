package controllers

import (
	"context"
	"errors"
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
		subscriber *Subscriber
		db         Database
	}

	CreateLedgerRequest struct {
		Identity domain.ID
		Name     string
	}

	CreateLedgerResponse struct {
		ID domain.ID
	}

	CreateExpenseRequest struct {
		Identity    domain.ID
		LedgerID    domain.ID
		Name        string
		ExpenseDate time.Time
		Records     []domain.NewRecord
	}

	CreateExpenseResponse struct {
		ID domain.ID
	}

	GetExpensesRequest struct {
		Identity domain.ID
		LedgerID domain.ID
		Cursor   time.Time
		Limit    int32
	}

	GetExpensesResponse struct {
		Expenses []v1.LedgerExpenseSummary
		Cursor   *time.Time
	}

	AddMembersRequest struct {
		Identity domain.ID
		LedgerID domain.ID
		Emails   []string
	}
)

func (c *Ledgers) Create(ctx context.Context, req CreateLedgerRequest) (resp *CreateLedgerResponse, err error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.Create")
	defer span.End()

	ctx = slog.Context(ctx, slog.WithStringer("identity", req.Identity))

	slog.Debug(ctx, "creating ledger", slog.WithAny("req", req))

	err = c.db.Transaction(ctx, func(db Database) error {
		err = db.Ledger().Create(ctx, req.Identity, func(count int64) (*domain.Ledger, error) {
			event, err := domain.NewLedger(req.Identity, req.Name, count+1)
			if err != nil {
				return nil, err
			}

			resp = &CreateLedgerResponse{
				ID: event.Data.ID,
			}

			return &event.Data, c.subscriber.handle(ctx, db, event)
		})
		if err != nil {
			return err
		}

		return nil
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
		slog.WithStringer("identity", req.Identity),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	event, err := domain.NewLedgerExpense(req.Identity, req.LedgerID, req.Name, req.ExpenseDate, req.Records)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "creating ledger expense", err)
	}
	resp = &CreateExpenseResponse{
		ID: event.Data.ID,
	}

	createFn := func(ledger *domain.Ledger) (*domain.Expense, error) {
		return &event.Data, nil
	}

	err = c.db.Transaction(ctx, func(db Database) error {
		if err := c.db.Expense().Create(ctx, req.LedgerID, createFn); err != nil {
			return err
		}
		return c.subscriber.handle(ctx, db, event)
	})
	switch {
	case errors.Is(err, domain.ErrUserNotAMember):
		if fieldErr := new(domain.FieldError); errors.As(err, fieldErr) {
			return nil, domain.FieldError{
				Cause:    domain.ErrUserNotAMember,
				Field:    fmt.Sprintf("user_balances.%d.identity", fieldErr.Metadata.Index),
				Metadata: fieldErr.Metadata,
			}
		}
		return nil, err
	case err != nil:
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
		slog.WithStringer("identity", req.Identity),
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

	ctx = slog.Context(ctx, slog.WithStringer("identity", identity))

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
func (c *Ledgers) AddParticipants(ctx context.Context, req AddMembersRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.AddParticipants")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("identity", req.Identity),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	transaction := func(db Database) error {
		users, err := db.User().ListByEmail(ctx, req.Emails)
		if err != nil {
			return slog.ErrorReturn(ctx, "failed to get users by email", err)
		}

		updateFn := func(ledger *domain.Ledger) error {
			pendingMemberIDs := kset.Select(func(u domain.User) domain.ID { return u.ID }, users...)
			events, err := ledger.AddParticipants(req.Identity, pendingMemberIDs...)
			if err != nil {
				return err
			}
			return c.subscriber.handle(ctx, db, convertEvents(events)...)
		}

		if err := db.Ledger().Update(ctx, req.LedgerID, updateFn); err != nil {
			return err
		}

		return nil
	}

	switch err := c.db.Transaction(ctx, transaction); {
	case err == nil:
		slog.Info(ctx, "added users to ledger")
		return nil
	case errors.Is(err, domain.ErrNotFound):
		return domain.FieldError{
			Field: "ledger_id",
			Cause: domain.ErrNotFound,
		}
	default:
		return slog.ErrorReturn(ctx, "failed to add users to ledger", err)
	}
}

func (c *Ledgers) GetParticipants(ctx context.Context, ledgerID domain.ID) ([]domain.LedgerParticipant, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetParticipants")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("ledger_id", ledgerID),
	)

	participants, err := c.db.Ledger().GetParticipants(ctx, ledgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger participants balances", err)
	}

	slog.Info(ctx, "ledger participants balances retrieved")

	return participants, nil
}
