package testfixtures

import (
	"maps"
	"slices"
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/kset"
	"github.com/stretchr/testify/require"
)

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

func Expense(t *testing.T, ledger *domain.Ledger, from, to domain.ID) *domain.Expense {
	currentMembers := kset.HashMapKey(slices.Collect(maps.Keys(ledger.Members))...)
	newMembers := kset.HashMapKey(from, to)

	err := ledger.AddMember(ledger.CreatedBy, newMembers.Difference(currentMembers).ToSlice()...)
	require.NoError(t, err)

	expense, err := ledger.CreateExpense(domain.CreateExpenseRequest{
		Creator:     ledger.CreatedBy,
		Name:        domain.NewID().String(),
		ExpenseDate: domain.Now(),
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
