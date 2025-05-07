package usecases

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

type (
	Ledgers interface {
		AddParticipants(ctx context.Context, req controllers.AddMembersRequest) error
		Create(ctx context.Context, req controllers.CreateRequest) (*controllers.CreateResponse, error)
		CreateExpense(ctx context.Context, req controllers.CreateExpenseRequest) (*controllers.CreateExpenseResponse, error)
		FindExpense(ctx context.Context, ledgerID v1.ID, expenseID v1.ID) (*v1.Expense, error)
		GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error)
		GetExpenses(ctx context.Context, req controllers.GetExpensesRequest) (*controllers.GetExpensesResult, error)
		GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error)
	}

	ledgers struct {
		ledgersController *controllers.Ledgers
		authorizer        Authorizer
	}
)

func NewLedgers(ledgersController *controllers.Ledgers) Ledgers {
	return &ledgers{
		ledgersController: ledgersController,
	}
}

func (l *ledgers) AddParticipants(ctx context.Context, req controllers.AddMembersRequest) error {
	if err := l.authorizer.Authorize(ctx, ActionLedgerExpenseWrite, ResourceUser(req.UserID), ResourceLedger(req.LedgerID)); err != nil {
		return err
	}

	return l.ledgersController.AddParticipants(ctx, req)
}

func (l *ledgers) Create(ctx context.Context, req controllers.CreateRequest) (*controllers.CreateResponse, error) {
	panic("unimplemented")
}

func (l *ledgers) CreateExpense(ctx context.Context, req controllers.CreateExpenseRequest) (*controllers.CreateExpenseResponse, error) {
	panic("unimplemented")
}

func (l *ledgers) FindExpense(ctx context.Context, ledgerID v1.ID, expenseID v1.ID) (*v1.Expense, error) {
	panic("unimplemented")
}

func (l *ledgers) GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error) {
	panic("unimplemented")
}

func (l *ledgers) GetExpenses(ctx context.Context, req controllers.GetExpensesRequest) (*controllers.GetExpensesResult, error) {
	panic("unimplemented")
}

func (l *ledgers) GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error) {
	panic("unimplemented")
}
