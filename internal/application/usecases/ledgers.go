package usecases

import (
	"context"
	"errors"

	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	Ledgers interface {
		AddParticipants(ctx context.Context, req controllers.AddMembersRequest) error
		Create(ctx context.Context, req controllers.CreateLedgerRequest) (*controllers.CreateLedgerResponse, error)
		CreateExpense(ctx context.Context, req controllers.CreateExpenseRequest) (*controllers.CreateExpenseResponse, error)
		FindExpense(ctx context.Context, identity domain.ID, ledgerID domain.ID, expenseID domain.ID) (*domain.Expense, error)
		GetByUser(ctx context.Context, identity domain.ID) ([]domain.Ledger, error)
		GetExpenses(ctx context.Context, req controllers.GetExpensesRequest) (*controllers.GetExpensesResponse, error)
		GetParticipants(ctx context.Context, identity domain.ID, ledgerID domain.ID) ([]domain.LedgerParticipant, error)
	}

	ledgers struct {
		controller *controllers.Controller
	}
)

func NewLedgers(controller *controllers.Controller) Ledgers {
	return &ledgers{
		controller: controller,
	}
}

func (l *ledgers) checkAuthorization(ctx context.Context, ledgerID domain.ID, f func(*domain.Ledger) bool) error {
	ledger, err := l.controller.Ledgers.Find(ctx, ledgerID)
	switch {
	case err == nil || errors.Is(err, domain.ErrNotFound):
		if err != nil || !f(ledger) {
			return domain.ErrForbidden
		}
		return nil
	default:
		return err
	}
}

func (l *ledgers) AddParticipants(ctx context.Context, req controllers.AddMembersRequest) error {
	if err := l.checkAuthorization(ctx, req.LedgerID, func(l *domain.Ledger) bool { return l.CreatedBy == req.Identity }); err != nil {
		return err
	}

	return l.controller.Ledgers.AddParticipants(ctx, req)
}

func (l *ledgers) Create(ctx context.Context, req controllers.CreateLedgerRequest) (*controllers.CreateLedgerResponse, error) {
	return l.controller.Ledgers.Create(ctx, req)
}

func (l *ledgers) CreateExpense(ctx context.Context, req controllers.CreateExpenseRequest) (*controllers.CreateExpenseResponse, error) {
	if err := l.checkAuthorization(ctx, req.LedgerID, func(l *domain.Ledger) bool { return l.IsParticipant(req.Identity) }); err != nil {
		return nil, err
	}

	return l.controller.Ledgers.CreateExpense(ctx, req)
}

func (l *ledgers) FindExpense(ctx context.Context, identity domain.ID, ledgerID domain.ID, expenseID domain.ID) (*domain.Expense, error) {
	if err := l.checkAuthorization(ctx, ledgerID, func(l *domain.Ledger) bool { return l.IsParticipant(identity) }); err != nil {
		return nil, err
	}

	return l.controller.Ledgers.FindExpense(ctx, ledgerID, expenseID)
}

func (l *ledgers) GetByUser(ctx context.Context, identity domain.ID) ([]domain.Ledger, error) {
	return l.controller.Ledgers.GetByIdentity(ctx, identity)
}

func (l *ledgers) GetExpenses(ctx context.Context, req controllers.GetExpensesRequest) (*controllers.GetExpensesResponse, error) {
	if err := l.checkAuthorization(ctx, req.LedgerID, func(l *domain.Ledger) bool { return l.IsParticipant(req.Identity) }); err != nil {
		return nil, err
	}

	return l.controller.Ledgers.GetExpenses(ctx, req)
}

func (l *ledgers) GetParticipants(ctx context.Context, identity domain.ID, ledgerID domain.ID) ([]domain.LedgerParticipant, error) {
	if err := l.checkAuthorization(ctx, ledgerID, func(l *domain.Ledger) bool { return l.IsParticipant(identity) }); err != nil {
		return nil, err
	}

	return l.controller.Ledgers.GetParticipants(ctx, ledgerID)
}
