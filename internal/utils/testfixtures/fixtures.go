package testfixtures

import (
	"testing"
	"time"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/require"
)

var Now = time.Now().UTC().Truncate(time.Microsecond)

func User(t *testing.T) *domain.User {
	user, err := domain.NewUser(domain.NewUserRequest{
		FirstName: domain.NewID().String(),
		LastName:  domain.NewID().String(),
		Email:     domain.NewID().String() + "@example.com",
		Password:  domain.NewID().String(),
	})
	require.NoError(t, err)

	return user
}

func Ledger(t *testing.T, creator *domain.User) *domain.Ledger {
	ledger, err := creator.CreateLedger(domain.NewID().String())
	require.NoError(t, err)
	return ledger
}

func Expense(t *testing.T, ledger *domain.Ledger) *domain.Expense {
	from := domain.NewID()
	to := domain.NewID()

	ledger.Members[from] = &domain.LedgerMember{}
	ledger.Members[to] = &domain.LedgerMember{}

	expense, err := ledger.CreateExpense(domain.CreateExpenseRequest{
		Creator:     ledger.CreatedBy,
		Name:        domain.NewID().String(),
		ExpenseDate: Now,
		PendingRecords: []domain.PendingRecord{
			{
				From:   from,
				To:     to,
				Type:   domain.RecordTypeDebt,
				Amount: 42,
			},
		},
	})
	require.NoError(t, err)
	return expense
}
