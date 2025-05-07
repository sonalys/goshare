package usecases

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	Ledgers interface {
		AddParticipants(ctx context.Context, req controllers.AddMembersRequest) error
		Create(ctx context.Context, req controllers.CreateLedgerRequest) (*controllers.CreateLedgerResponse, error)
		CreateExpense(ctx context.Context, req controllers.CreateExpenseRequest) (*controllers.CreateExpenseResponse, error)
		FindExpense(ctx context.Context, ledgerID domain.ID, expenseID domain.ID) (*domain.Expense, error)
		GetByUser(ctx context.Context, userID domain.ID) ([]domain.Ledger, error)
		GetExpenses(ctx context.Context, req controllers.GetExpensesRequest) (*controllers.GetExpensesResponse, error)
		GetParticipants(ctx context.Context, ledgerID domain.ID) ([]domain.LedgerParticipant, error)
	}

	ledgers struct {
		ledgersController *controllers.Ledgers
	}
)

func NewLedgers(ledgersController *controllers.Ledgers) Ledgers {
	return &ledgers{
		ledgersController: ledgersController,
	}
}

func (l *ledgers) AddParticipants(ctx context.Context, req controllers.AddMembersRequest) error {
	return l.ledgersController.AddParticipants(ctx, req)
}

func (l *ledgers) Create(ctx context.Context, req controllers.CreateLedgerRequest) (*controllers.CreateLedgerResponse, error) {
	return l.ledgersController.Create(ctx, req)
}

func (l *ledgers) CreateExpense(ctx context.Context, req controllers.CreateExpenseRequest) (*controllers.CreateExpenseResponse, error) {
	return l.ledgersController.CreateExpense(ctx, req)
}

func (l *ledgers) FindExpense(ctx context.Context, ledgerID domain.ID, expenseID domain.ID) (*domain.Expense, error) {
	return l.ledgersController.FindExpense(ctx, ledgerID, expenseID)
}

func (l *ledgers) GetByUser(ctx context.Context, userID domain.ID) ([]domain.Ledger, error) {
	return l.ledgersController.GetByUser(ctx, userID)
}

func (l *ledgers) GetExpenses(ctx context.Context, req controllers.GetExpensesRequest) (*controllers.GetExpensesResponse, error) {
	return l.ledgersController.GetExpenses(ctx, req)
}

func (l *ledgers) GetParticipants(ctx context.Context, ledgerID domain.ID) ([]domain.LedgerParticipant, error) {
	return l.ledgersController.GetParticipants(ctx, ledgerID)
}
