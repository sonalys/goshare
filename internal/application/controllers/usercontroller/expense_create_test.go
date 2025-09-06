package usercontroller_test

import (
	"context"
	"testing"

	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Expense_Create(t *testing.T) {
	t.Parallel()

	type testSetup struct {
		db *databaseMock
	}

	createTestData := func() usercontroller.CreateExpenseRequest {
		return usercontroller.CreateExpenseRequest{
			ActorID:     domain.NewID(),
			LedgerID:    domain.NewID(),
			Name:        domain.NewID().String(),
			ExpenseDate: domain.Now(),
			PendingRecords: []domain.PendingRecord{
				{
					From:   domain.NewID(),
					To:     domain.NewID(),
					Type:   domain.RecordTypeDebt,
					Amount: 42,
				},
			},
		}
	}

	setup := func(t *testing.T, td usercontroller.CreateExpenseRequest) (usercontroller.Controller, testSetup) {
		mocks := testSetup{
			db: setupDatabaseMock(t),
		}

		user := testfixtures.User(t)
		user.ID = td.ActorID

		ledger := testfixtures.Ledger(t, user)
		ledger.ID = td.LedgerID

		for _, records := range td.PendingRecords {
			ledger.Members[records.From] = &domain.LedgerMember{}
			ledger.Members[records.To] = &domain.LedgerMember{}
		}

		mocks.db.tx.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			assert.Equal(t, td.LedgerID, id)

			return ledger, nil
		}

		mocks.db.tx.expense.CreateFunc = func(ctx context.Context, expense *domain.Expense) error {
			assert.Equal(t, td.LedgerID, expense.LedgerID)

			return nil
		}

		mocks.db.tx.ledger.UpdateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			assert.Equal(t, td.LedgerID, ledger.ID)

			return nil
		}

		controller := usercontroller.New(usercontroller.Dependencies{
			LocalDatabase: &mocks.db.db,
		})

		return controller, mocks
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)

		resp, err := controller.Expenses().Create(ctx, td)
		require.NoError(t, err)
		require.NotEmpty(t, resp)
	})

	t.Run("fail/user is not authorized", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, _ := setup(t, td)
		td.ActorID = domain.NewID()

		resp, err := controller.Expenses().Create(ctx, td)
		require.ErrorIs(t, err, application.ErrForbidden)
		require.Empty(t, resp)
	})

	t.Run("fail/ledger repository get errors", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.GetFunc = func(ctx context.Context, id domain.ID) (*domain.Ledger, error) {
			return nil, assert.AnError
		}

		resp, err := controller.Expenses().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Empty(t, resp)
	})

	t.Run("fail/expense repository create errors", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.expense.CreateFunc = func(ctx context.Context, expense *domain.Expense) error {
			return assert.AnError
		}

		resp, err := controller.Expenses().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Empty(t, resp)
	})

	t.Run("fail/ledger repository update errors", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		td := createTestData()
		controller, mocks := setup(t, td)

		mocks.db.tx.ledger.UpdateFunc = func(ctx context.Context, ledger *domain.Ledger) error {
			return assert.AnError
		}

		resp, err := controller.Expenses().Create(ctx, td)
		require.ErrorIs(t, err, assert.AnError)
		require.Empty(t, resp)
	})
}
