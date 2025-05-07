package controllers

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/sonalys/goshare/internal/application/pkg/otel"
	"github.com/sonalys/goshare/internal/application/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

type (
	Ledgers struct {
		subscriber *Subscriber
		db         Database
	}
)

type (
	CreateRequest struct {
		UserID v1.ID
		Name   string
	}

	CreateResponse struct {
		ID v1.ID
	}
)

func (r CreateRequest) Validate() error {
	var errs v1.FormError

	if r.UserID.IsEmpty() {
		errs = append(errs, v1.NewRequiredFieldError("user_id"))
	}

	if r.Name == "" {
		errs = append(errs, v1.NewRequiredFieldError("name"))
	} else if nameLength := len(r.Name); nameLength < 3 || nameLength > 255 {
		errs = append(errs, v1.NewFieldLengthError("name", 3, 255))
	}

	return errs.Validate()
}

func (c *Ledgers) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.Create")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("user_id", req.UserID),
	)

	slog.Debug(ctx, "creating ledger", slog.WithAny("req", req))

	if err := req.Validate(); err != nil {
		return nil, slog.ErrorReturn(ctx, "invalid request", err)
	}

	ledger := &v1.Ledger{
		ID:   v1.NewID(),
		Name: req.Name,
		Participants: []v1.LedgerParticipant{
			{
				ID:        v1.NewID(),
				UserID:    req.UserID,
				Balance:   0,
				CreatedAt: time.Now(),
				CreatedBy: req.UserID,
			},
		},
		CreatedAt: time.Now(),
		CreatedBy: req.UserID,
	}

	ctx = slog.Context(ctx,
		slog.WithStringer("ledger_id", ledger.ID),
	)

	slog.Debug(ctx, "ledger entity initialized", slog.WithAny("ledger", ledger))

	err := c.db.Ledger().Create(ctx, req.UserID, func(count int64) (*v1.Ledger, error) {
		if count+1 > v1.UserMaxLedgers {
			return nil, v1.ErrUserMaxLedgers
		}

		return ledger, nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to create ledger", err)
	}

	slog.Info(ctx, "ledger created")

	resp := &CreateResponse{
		ID: ledger.ID,
	}

	return resp, nil
}

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

	GetExpensesRequest struct {
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

func (c *Ledgers) CreateExpense(ctx context.Context, req CreateExpenseRequest) (*CreateExpenseResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.CreateExpense")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("user_id", req.UserID),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	if err := req.Validate(); err != nil {
		return nil, slog.ErrorReturn(ctx, "validating request", err)
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

	ctx = slog.Context(ctx,
		slog.WithStringer("expense_id", expense.ID),
	)

	switch err := c.db.Expense().Create(ctx, req.LedgerID, func(ledger *v1.Ledger) (*v1.Expense, error) {
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
		return nil, slog.ErrorReturn(ctx, "creating expense", err)
	default:
		slog.Info(ctx, "expense created")

		return &CreateExpenseResponse{
			ID: expense.ID,
		}, nil
	}
}

func (c *Ledgers) FindExpense(ctx context.Context, ledgerID v1.ID, expenseID v1.ID) (*v1.Expense, error) {
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

func (c *Ledgers) GetExpenses(ctx context.Context, req GetExpensesRequest) (*GetExpensesResult, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetExpenses")
	defer span.End()

	req.Limit = max(1, req.Limit)

	ctx = slog.Context(ctx,
		slog.WithStringer("user_id", req.UserID),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	ledger, err := c.db.Ledger().Find(ctx, req.LedgerID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger", err)
	}

	if !ledger.IsParticipant(req.UserID) {
		return nil, slog.ErrorReturn(ctx, "authorizing request", v1.ErrUserNotAMember)
	}

	expenses, err := c.db.Expense().GetByLedger(ctx, req.LedgerID, req.Cursor, req.Limit+1)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to get ledger expenses", err)
	}

	if len(expenses) == 0 {
		return &GetExpensesResult{}, nil
	}

	slog.Info(ctx, "ledger expenses retrieved")

	var cursor *time.Time
	if len(expenses) == int(req.Limit)+1 {
		expenses = expenses[:len(expenses)-1]
		cursor = &expenses[len(expenses)-1].CreatedAt
	}

	return &GetExpensesResult{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}

func (c *Ledgers) GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.ListByUser")
	defer span.End()

	ledgers, err := c.db.Ledger().GetByUser(ctx, userID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to list ledgers", err)
	}

	return ledgers, nil
}

type (
	AddMembersRequest struct {
		UserID   v1.ID
		LedgerID v1.ID
		Emails   []string
	}
)

func (r *AddMembersRequest) Validate() error {
	var errs v1.FormError

	if r.UserID.IsEmpty() {
		errs = append(errs, v1.NewRequiredFieldError("user_id"))
	}

	if r.LedgerID.IsEmpty() {
		errs = append(errs, v1.NewRequiredFieldError("ledger_id"))
	}

	r.Emails = slices.Compact(r.Emails)
	switch lenEmails := len(r.Emails); {
	case lenEmails == 0:
		errs = append(errs, v1.NewRequiredFieldError("emails"))
	case lenEmails > v1.LedgerMaxUsers-1:
		errs = append(errs, v1.NewFieldLengthError("emails", 1, v1.LedgerMaxUsers))
	}

	return errs.Validate()
}

// TODO(invitations): Here it's a simplification of the user membership process.
// We can always invert the flow and create invitation links, so the users click themselves
// We can also send invites through the system and they accept the invite through the API.
func (c *Ledgers) AddParticipants(ctx context.Context, req AddMembersRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.AddMembers")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("user_id", req.UserID),
		slog.WithStringer("ledger_id", req.LedgerID),
	)

	users, err := c.db.User().ListByEmail(ctx, req.Emails)
	if err != nil {
		return slog.ErrorReturn(ctx, "failed to get users by email", err)
	}

	ids := make([]v1.ID, 0, len(users))
	for _, user := range users {
		if user.ID == req.UserID {
			continue
		}
		ids = append(ids, user.ID)
	}

	err = c.db.Ledger().AddParticipants(ctx, req.LedgerID, func(ledger *v1.Ledger) error {
		if ledger.CreatedBy != req.UserID {
			return fmt.Errorf("user %s is not the owner of the ledger %s", req.UserID, ledger.ID)
		}

		ledger.AddParticipants(req.UserID, ids...)

		if len(ledger.Participants) >= v1.LedgerMaxUsers {
			return v1.ErrLedgerMaxUsers
		}

		return nil
	})
	switch {
	case err == nil:
		slog.Info(ctx, "added users to ledger")
		return nil
	case errors.Is(err, v1.ErrNotFound):
		return v1.FieldError{
			Field: "ledger_id",
			Cause: v1.ErrNotFound,
		}
	default:
		return slog.ErrorReturn(ctx, "failed to add users to ledger", err)
	}
}

func (c *Ledgers) GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetBalances")
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
